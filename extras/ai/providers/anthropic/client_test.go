package anthropic

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
		if r.URL.Path != "/messages" {
			t.Errorf("path = %q, want /messages", r.URL.Path)
		}
		if got := r.Header.Get("x-api-key"); got != "sk-ant-test" {
			t.Errorf("x-api-key = %q", got)
		}
		if got := r.Header.Get("anthropic-version"); got != APIVersion {
			t.Errorf("anthropic-version = %q", got)
		}
		var body map[string]any
		b, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(b, &body)
		if body["model"] != DefaultModel {
			t.Errorf("model = %v", body["model"])
		}
		if body["max_tokens"] != 1024.0 {
			t.Errorf("max_tokens = %v", body["max_tokens"])
		}
		if body["system"] != "be brief" {
			t.Errorf("system = %v", body["system"])
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id":    "msg_1",
			"model": body["model"],
			"content": []map[string]any{
				{"type": "text", "text": "hello from claude"},
			},
			"usage": map[string]any{"input_tokens": 5, "output_tokens": 3},
		})
	}))
	defer srv.Close()

	c := New(WithAPIKey("sk-ant-test"), WithBaseURL(srv.URL))
	res, err := c.Generate(context.Background(), ai.Request{Prompt: "hi", System: "be brief"})
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if res.Text != "hello from claude" {
		t.Fatalf("Text = %q", res.Text)
	}
	if res.Provider != Name {
		t.Fatalf("Provider = %q", res.Provider)
	}
	if res.Usage == nil || res.Usage.TotalTokens != 8 {
		t.Fatalf("Usage = %+v", res.Usage)
	}
}

func TestGenerateReasoningEnablesThinking(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		b, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(b, &body)
		th, _ := body["thinking"].(map[string]any)
		if th == nil {
			t.Fatal("expected thinking block")
		}
		if th["type"] != "enabled" {
			t.Errorf("thinking.type = %v", th["type"])
		}
		if th["budget_tokens"] != 32768.0 {
			t.Errorf("budget = %v", th["budget_tokens"])
		}
		// max_tokens must exceed the budget.
		if body["max_tokens"].(float64) <= th["budget_tokens"].(float64) {
			t.Errorf("max_tokens not greater than budget")
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"content": []map[string]any{{"type": "text", "text": "ok"}},
		})
	}))
	defer srv.Close()

	c := New(WithAPIKey("k"), WithBaseURL(srv.URL))
	_, err := c.Generate(context.Background(), ai.Request{Prompt: "hi", Reasoning: ai.ReasoningHigh})
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
}

func TestGenerateErrors(t *testing.T) {
	if _, err := New().Generate(context.Background(), ai.Request{Prompt: "hi"}); err == nil {
		t.Fatal("expected missing key error")
	}
	c := New(WithAPIKey("k"))
	if _, err := c.Generate(context.Background(), ai.Request{}); err == nil {
		t.Fatal("expected missing prompt error")
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()
	if _, err := New(WithAPIKey("k"), WithBaseURL(srv.URL)).Generate(context.Background(), ai.Request{Prompt: "hi"}); err == nil {
		t.Fatal("expected HTTP error")
	}
}
