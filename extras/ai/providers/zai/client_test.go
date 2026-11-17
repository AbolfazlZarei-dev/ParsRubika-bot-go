package zai

import (
	"context"
	"encoding/json"
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

func TestGenerate(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chat/completions" {
			t.Errorf("path = %q", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer zai-key" {
			t.Errorf("Authorization = %q", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"choices": []map[string]any{
				{"message": map[string]any{"role": "assistant", "content": "glm reply"}},
			},
		})
	}))
	defer srv.Close()

	c := New(WithAPIKey("zai-key"), WithBaseURL(srv.URL))
	res, err := c.Generate(context.Background(), ai.Request{Prompt: "hi"})
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if res.Text != "glm reply" {
		t.Fatalf("Text = %q", res.Text)
	}
	if res.Provider != Name {
		t.Fatalf("Provider = %q", res.Provider)
	}
}
