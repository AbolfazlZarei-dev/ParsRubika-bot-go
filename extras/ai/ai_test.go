package ai

import (
	"context"
	"errors"
	"strings"
	"testing"
)

// fakeProvider is a minimal Provider used to exercise the adapter and
// registry without any network.
type fakeProvider struct {
	name    string
	model   string
	lastReq Request
	reply   string
	err     error
}

func (f *fakeProvider) Name() string         { return f.name }
func (f *fakeProvider) DefaultModel() string { return f.model }
func (f *fakeProvider) Generate(ctx context.Context, req Request) (Response, error) {
	f.lastReq = req
	if f.err != nil {
		return Response{}, f.err
	}
	model := req.Model
	if model == "" {
		model = f.model
	}
	return Response{Text: f.reply, Model: model, Provider: f.name}, nil
}

func TestExtractPrompt(t *testing.T) {
	cases := []struct {
		text, cmd string
		prefixes  []string
		want      string
	}{
		{"/ask hello", "ask", []string{"/"}, "hello"},
		{"  /Ask   world  ", "ask", []string{"/"}, "world"},
		{"ask foo bar", "ask", []string{"/"}, "foo bar"},
		{"/ask", "ask", []string{"/"}, ""},
		{"", "ask", []string{"/"}, ""},
		{"!ask hi", "ask", []string{"/", "!"}, "hi"},
		{"hello", "ask", []string{"/"}, ""}, // not a command
	}
	for _, c := range cases {
		got := extractPrompt(c.text, c.cmd, c.prefixes)
		if got != c.want {
			t.Errorf("extractPrompt(%q, %q, %v) = %q, want %q", c.text, c.cmd, c.prefixes, got, c.want)
		}
	}
}

func TestNewGeneratorAdapter(t *testing.T) {
	p := &fakeProvider{name: "fake", model: "fake-1", reply: "pong"}
	g := NewGenerator(p, WithGeneratorSystem("be nice"), WithGeneratorModel("fake-2"))

	out, err := g(context.Background(), "ping")
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if out != "pong" {
		t.Fatalf("output = %q, want pong", out)
	}
	if p.lastReq.System != "be nice" {
		t.Errorf("System not forwarded: %q", p.lastReq.System)
	}
	if p.lastReq.Model != "fake-2" {
		t.Errorf("Model not forwarded: %q", p.lastReq.Model)
	}
	if p.lastReq.Prompt != "ping" {
		t.Errorf("Prompt not forwarded: %q", p.lastReq.Prompt)
	}
}

func TestNewGeneratorDefaultModel(t *testing.T) {
	p := &fakeProvider{name: "fake", model: "default-model", reply: "ok"}
	g := NewGenerator(p) // no model override
	if _, err := g(context.Background(), "hi"); err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if p.lastReq.Model != "" {
		t.Errorf("empty Request.Model expected, got %q", p.lastReq.Model)
	}
}

func TestNewGeneratorNilProvider(t *testing.T) {
	g := NewGenerator(nil)
	if _, err := g(context.Background(), "hi"); err == nil {
		t.Fatal("expected error from nil-adapted generator")
	}
}

func TestNewGeneratorForwardsError(t *testing.T) {
	want := errors.New("boom")
	p := &fakeProvider{name: "fake", model: "m", err: want}
	g := NewGenerator(p)
	if _, err := g(context.Background(), "hi"); !errors.Is(err, want) {
		t.Fatalf("expected boom, got %v", err)
	}
}

func TestRegisterProvider(t *testing.T) {
	// Snapshot and restore the global registry so parallel tests do
	// not stomp on each other.
	registryMu.Lock()
	prev := registry
	registry = nil
	registryMu.Unlock()
	defer func() {
		registryMu.Lock()
		registry = prev
		registryMu.Unlock()
	}()

	if g := registeredGenerator(); g != nil {
		t.Fatal("expected nil registry at start")
	}
	RegisterProvider(&fakeProvider{name: "fake", model: "m", reply: "registered"})
	g := registeredGenerator()
	if g == nil {
		t.Fatal("expected generator after RegisterProvider")
	}
	out, err := g(context.Background(), "hi")
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if !strings.Contains(out, "registered") {
		t.Fatalf("output = %q, want 'registered'", out)
	}
}

// TestProviderInterfaceConformance is a compile-time check that
// fakeProvider satisfies the Provider interface.
func TestProviderInterfaceConformance(t *testing.T) {
	var _ Provider = (*fakeProvider)(nil)
}
