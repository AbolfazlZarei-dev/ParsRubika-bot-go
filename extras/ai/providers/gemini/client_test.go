package gemini

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
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
	var gotURL string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotURL = r.URL.String()
		var body map[string]any
		b, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(b, &body)
		contents, _ := body["contents"].([]any)
		if len(contents) != 1 {
			t.Errorf("contents len = %d", len(contents))
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"candidates": []map[string]any{
				{"content": map[string]any{"parts": []map[string]any{{"text": "hello from gemini"}}}, "finishReason": "STOP"},
			},
			"usageMetadata": map[string]any{"promptTokenCount": 4, "candidatesTokenCount": 3, "totalTokenCount": 7},
		})
	}))
	defer srv.Close()

	c := New(WithAPIKey("AIza-test"), WithBaseURL(srv.URL))
	res, err := c.Generate(context.Background(), ai.Request{Prompt: "hi", System: "be brief"})
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if res.Text != "hello from gemini" {
		t.Fatalf("Text = %q", res.Text)
	}
	if !strings.Contains(gotURL, "/models/"+DefaultModel+":generateContent") {
		t.Errorf("URL = %q, missing model segment", gotURL)
	}
	if !strings.Contains(gotURL, "key=AIza-test") {
		t.Errorf("URL = %q, missing key param", gotURL)
	}
	if res.Usage == nil || res.Usage.TotalTokens != 7 {
		t.Fatalf("Usage = %+v", res.Usage)
	}
}

func TestGenerateModelInURL(t *testing.T) {
	var gotURL string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotURL = r.URL.String()
		_ = json.NewEncoder(w).Encode(map[string]any{
			"candidates": []map[string]any{{"content": map[string]any{"parts": []map[string]any{{"text": "ok"}}}}},
		})
	}))
	defer srv.Close()

	c := New(WithAPIKey("k"), WithBaseURL(srv.URL))
	_, err := c.Generate(context.Background(), ai.Request{Model: "gemini-3.8-flash", Prompt: "hi"})
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if !strings.Contains(gotURL, "/models/gemini-3.8-flash:generateContent") {
		t.Errorf("URL = %q, expected model in path", gotURL)
	}
}

func TestGenerateErrors(t *testing.T) {
	if _, err := New().Generate(context.Background(), ai.Request{Prompt: "hi"}); err == nil {
		t.Fatal("expected missing key error")
	}
	if _, err := New(WithAPIKey("k")).Generate(context.Background(), ai.Request{}); err == nil {
		t.Fatal("expected missing prompt error")
	}
}
