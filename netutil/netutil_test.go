package netutil

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestValidateURLScheme(t *testing.T) {
	cases := []struct {
		url   string
		block bool
	}{
		{"http://example.com/x", false},
		{"https://example.com/x", false},
		{"file:///etc/passwd", true},
		{"ftp://example.com/x", true},
		{"gopher://127.0.0.1:6379/_FLUSHALL", true},
		{"", true},
	}
	for _, c := range cases {
		_, err := ValidateURL(c.url)
		if c.block && err == nil {
			t.Errorf("ValidateURL(%q): want block, got nil", c.url)
		}
		if !c.block && err != nil {
			t.Errorf("ValidateURL(%q): want ok, got %v", c.url, err)
		}
	}
}

func TestValidateURLPrivateIPs(t *testing.T) {
	for _, u := range []string{
		"http://127.0.0.1/x",
		"http://10.0.0.1/x",
		"http://192.168.1.1/x",
		"http://172.16.0.1/x",
		"http://169.254.169.254/latest/meta-data/", // cloud metadata
		"http://localhost/x",
		"http://api.localhost/x",
		"http://[::1]/x",
		"http://[fe80::1]/x",
	} {
		if _, err := ValidateURL(u); err == nil {
			t.Errorf("ValidateURL(%q): want block, got nil", u)
		}
	}
}

func TestAssertURLResolvesPrivateBlocked(t *testing.T) {
	// "localhost" resolves to 127.0.0.1 — must be blocked at resolve time.
	if _, err := AssertURL(context.Background(), "http://localhost/"); err == nil {
		t.Fatal("AssertURL(localhost) must be blocked")
	}
	// A public hostname must pass (we use example.com which is reserved
	// for documentation and resolves to public IPs).
	if _, err := AssertURL(context.Background(), "https://example.com/"); err != nil {
		t.Fatalf("AssertURL(example.com): %v", err)
	}
}

func TestSafeTransportBlocksLoopbackDial(t *testing.T) {
	// Spin a server on the loopback and confirm a SafeClient cannot reach it.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("secret"))
	}))
	defer srv.Close()
	// httptest.Server listens on 127.0.0.1 — ValidateURL already blocks it,
	// but bypass the pre-flight to exercise the dial-layer guard directly.
	client := NewSafeClient(0)
	resp, err := client.Get(srv.URL)
	if err == nil {
		resp.Body.Close()
		t.Fatalf("SafeClient must refuse dial to loopback server, got status %d", resp.StatusCode)
	}
}

func TestSafeTransportAllowsPublic(t *testing.T) {
	// We can only safely assert the transport construction + a public URL
	// pre-flight without making a real network round-trip; the dial-layer
	// guard is exercised by the loopback test above. Here we just confirm
	// a non-blocked URL passes ValidateURL.
	if _, err := ValidateURL("https://rubika.ir/post/abc"); err != nil {
		t.Fatalf("ValidateURL(public): %v", err)
	}
}
