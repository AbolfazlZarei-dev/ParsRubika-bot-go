package openrouter

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/extras/ai"
)

func TestNewDefaults(t *testing.T) {
	c := New()
	if c.Name() != Name || c.DefaultModel() != DefaultModel || c.baseURL != DefaultBaseURL {
		t.Errorf("defaults wrong: %+v", c)
	}
}

func TestGenerateWithRankings(t *testing.T) {
	var gotReferer, gotTitle, gotModel string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chat/completions" {
			t.Errorf("path = %q", r.URL.Path)
		}
		gotReferer = r.Header.Get("HTTP-Referer")
		gotTitle = r.Header.Get("X-Title")
		if got := r.Header.Get("Authorization"); got != "Bearer sk-or-test" {
			t.Errorf("Authorization = %q", got)
		}
		var body map[string]any
		b, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(b, &body)
		gotModel, _ = body["model"].(string)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"choices": []map[string]any{
				{"message": map[string]any{"role": "assistant", "content": "routed reply"}},
			},
		})
	}))
	defer srv.Close()

	c := New(
		WithAPIKey("sk-or-test"),
		WithBaseURL(srv.URL),
		WithRankings("https://mybot.example", "My Rubika Bot"),
	)
	res, err := c.Generate(context.Background(), ai.Request{
		Model:  "anthropic/claude-opus-5-5",
		Prompt: "hi",
	})
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if res.Text != "routed reply" {
		t.Fatalf("Text = %q", res.Text)
	}
	if gotModel != "anthropic/claude-opus-5-5" {
		t.Errorf("model forwarded = %q", gotModel)
	}
	if gotReferer != "https://mybot.example" {
		t.Errorf("HTTP-Referer = %q", gotReferer)
	}
	if gotTitle != "My Rubika Bot" {
		t.Errorf("X-Title = %q", gotTitle)
	}
}
