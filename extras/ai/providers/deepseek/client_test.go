package deepseek

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

func TestGenerate(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chat/completions" {
			t.Errorf("path = %q", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer sk-ds" {
			t.Errorf("Authorization = %q", got)
		}
		var body map[string]any
		b, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(b, &body)
		if body["model"] != DefaultModel {
			t.Errorf("model = %v", body["model"])
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"choices": []map[string]any{
				{"message": map[string]any{"role": "assistant", "content": "deepseek says hi"}},
			},
		})
	}))
	defer srv.Close()

	c := New(WithAPIKey("sk-ds"), WithBaseURL(srv.URL))
	res, err := c.Generate(context.Background(), ai.Request{Prompt: "hi"})
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if res.Text != "deepseek says hi" {
		t.Fatalf("Text = %q", res.Text)
	}
	if res.Provider != Name {
		t.Fatalf("Provider = %q", res.Provider)
	}
}
