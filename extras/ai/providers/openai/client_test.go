package openai

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/extras/ai"
)

// newChatServer stands up a fake OpenAI-style /chat/completions endpoint
// and asserts the bearer token + model wiring.
func newChatServer(t *testing.T, wantModel, wantKey string, reply string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chat/completions" {
			t.Errorf("path = %q, want /chat/completions", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer "+wantKey {
			t.Errorf("Authorization = %q, want Bearer %s", got, wantKey)
		}
		var body map[string]any
		b, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(b, &body)
		if body["model"] != wantModel {
			t.Errorf("model = %v, want %v", body["model"], wantModel)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"model": body["model"],
			"choices": []map[string]any{
				{"index": 0, "message": map[string]any{"role": "assistant", "content": reply}, "finish_reason": "stop"},
			},
		})
	}))
}

func TestNewDefaults(t *testing.T) {
	c := New()
	if c.Name() != Name {
		t.Errorf("Name = %q", c.Name())
	}
	if c.DefaultModel() != DefaultModel {
		t.Errorf("DefaultModel = %q", c.DefaultModel())
	}
	if c.baseURL != DefaultBaseURL {
		t.Errorf("baseURL = %q", c.baseURL)
	}
}

func TestGenerate(t *testing.T) {
	srv := newChatServer(t, DefaultModel, "sk-test", "hello")
	defer srv.Close()

	c := New(WithAPIKey("sk-test"), WithBaseURL(srv.URL))
	res, err := c.Generate(context.Background(), ai.Request{Prompt: "hi"})
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if res.Text != "hello" {
		t.Fatalf("Text = %q", res.Text)
	}
	if res.Provider != Name {
		t.Fatalf("Provider = %q", res.Provider)
	}
}

func TestGenerateModelOverride(t *testing.T) {
	srv := newChatServer(t, "gpt-4o", "sk-test", "ok")
	defer srv.Close()

	c := New(WithAPIKey("sk-test"), WithBaseURL(srv.URL), WithModel("gpt-4o"))
	res, err := c.Generate(context.Background(), ai.Request{Prompt: "hi"})
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if res.Text != "ok" {
		t.Fatalf("Text = %q", res.Text)
	}
}

func TestGeneratorAdapter(t *testing.T) {
	srv := newChatServer(t, DefaultModel, "sk-test", "pong")
	defer srv.Close()

	g := Generator(WithAPIKey("sk-test"), WithBaseURL(srv.URL))
	out, err := g(context.Background(), "ping")
	if err != nil {
		t.Fatalf("Generator: %v", err)
	}
	if out != "pong" {
		t.Fatalf("output = %q", out)
	}
}
