package gapgpt

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

func TestGenerateResponsesShape(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/responses" {
			t.Errorf("path = %q, want /responses", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer gap-test" {
			t.Errorf("Authorization = %q", got)
		}
		var body map[string]any
		b, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(b, &body)
		if body["model"] != DefaultModel {
			t.Errorf("model = %v", body["model"])
		}
		if body["input"] != "hi" {
			t.Errorf("input = %v", body["input"])
		}
		if rsn, _ := body["reasoning"].(map[string]any); rsn == nil || rsn["effort"] != "low" {
			t.Errorf("reasoning = %v", body["reasoning"])
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id":    "resp_1",
			"model": body["model"],
			"output": []map[string]any{
				{"type": "message", "content": []map[string]any{{"type": "output_text", "text": "once upon a time"}}},
			},
			"usage": map[string]any{"input_tokens": 3, "output_tokens": 5, "total_tokens": 8},
		})
	}))
	defer srv.Close()

	c := New(WithAPIKey("gap-test"), WithBaseURL(srv.URL))
	res, err := c.Generate(context.Background(), ai.Request{Prompt: "hi", Reasoning: ai.ReasoningLow})
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if res.Text != "once upon a time" {
		t.Fatalf("Text = %q", res.Text)
	}
	if res.Provider != Name {
		t.Fatalf("Provider = %q", res.Provider)
	}
	if res.Usage == nil || res.Usage.TotalTokens != 8 {
		t.Fatalf("Usage = %+v", res.Usage)
	}
}

func TestGenerateChatCompletionsFallback(t *testing.T) {
	// Some gateways answer /responses with a chat-completions body; the
	// extractor must still find the text.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"choices": []map[string]any{
				{"message": map[string]any{"role": "assistant", "content": "fallback ok"}},
			},
		})
	}))
	defer srv.Close()

	c := New(WithAPIKey("k"), WithBaseURL(srv.URL))
	res, err := c.Generate(context.Background(), ai.Request{Prompt: "hi"})
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if res.Text != "fallback ok" {
		t.Fatalf("Text = %q", res.Text)
	}
}

func TestGenerateSystemFoldedIntoInput(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		b, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(b, &body)
		input, _ := body["input"].(string)
		if input[:9] != "be brief\n" {
			t.Errorf("system not folded into input: %q", input)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"output_text": "ok"})
	}))
	defer srv.Close()

	c := New(WithAPIKey("k"), WithBaseURL(srv.URL))
	_, err := c.Generate(context.Background(), ai.Request{Prompt: "hi", System: "be brief"})
	if err != nil {
		t.Fatalf("Generate: %v", err)
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
