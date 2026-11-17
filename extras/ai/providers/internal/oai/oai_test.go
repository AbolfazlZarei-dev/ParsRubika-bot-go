package oai

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/extras/ai"
)

func newServer(t *testing.T, wantBody func(map[string]any), reply string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chat/completions" {
			t.Errorf("path = %q, want /chat/completions", r.URL.Path)
		}
		if got, want := r.Header.Get("Authorization"), "Bearer sk-test"; got != want {
			t.Errorf("Authorization = %q, want %q", got, want)
		}
		b, _ := io.ReadAll(r.Body)
		var m map[string]any
		_ = json.Unmarshal(b, &m)
		wantBody(m)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id":    "chatcmpl-1",
			"model": m["model"],
			"choices": []map[string]any{
				{"index": 0, "message": map[string]any{"role": "assistant", "content": reply}, "finish_reason": "stop"},
			},
			"usage": map[string]any{"prompt_tokens": 4, "completion_tokens": 2, "total_tokens": 6},
		})
	}))
}

func TestChatHappyPath(t *testing.T) {
	srv := newServer(t, func(m map[string]any) {
		if m["model"] != "gpt-x" {
			t.Errorf("model = %v, want gpt-x", m["model"])
		}
		msgs, _ := m["messages"].([]any)
		if len(msgs) != 1 {
			t.Errorf("messages len = %d, want 1", len(msgs))
		}
	}, "hello back")
	defer srv.Close()

	cfg := Config{BaseURL: srv.URL, APIKey: "sk-test", DefaultModel: "fallback"}
	res, err := Chat(context.Background(), cfg, "openai", ai.Request{Model: "gpt-x", Prompt: "hi"})
	if err != nil {
		t.Fatalf("Chat: %v", err)
	}
	if res.Text != "hello back" {
		t.Fatalf("Text = %q", res.Text)
	}
	if res.Model != "gpt-x" {
		t.Fatalf("Model = %q", res.Model)
	}
	if res.Provider != "openai" {
		t.Fatalf("Provider = %q", res.Provider)
	}
	if res.Usage == nil || res.Usage.TotalTokens != 6 {
		t.Fatalf("Usage = %+v", res.Usage)
	}
}

func TestChatDefaultModel(t *testing.T) {
	srv := newServer(t, func(m map[string]any) {
		if m["model"] != "fallback" {
			t.Errorf("model = %v, want fallback", m["model"])
		}
	}, "ok")
	defer srv.Close()

	cfg := Config{BaseURL: srv.URL, APIKey: "sk-test", DefaultModel: "fallback"}
	res, err := Chat(context.Background(), cfg, "openai", ai.Request{Prompt: "hi"})
	if err != nil {
		t.Fatalf("Chat: %v", err)
	}
	if res.Text != "ok" {
		t.Fatalf("Text = %q", res.Text)
	}
}

func TestChatSystemAndParams(t *testing.T) {
	srv := newServer(t, func(m map[string]any) {
		msgs, _ := m["messages"].([]any)
		if len(msgs) != 2 {
			t.Fatalf("messages len = %d, want 2 (system+user)", len(msgs))
		}
		first, _ := msgs[0].(map[string]any)
		if first["role"] != "system" {
			t.Errorf("first role = %v, want system", first["role"])
		}
		if m["temperature"] != 0.7 {
			t.Errorf("temperature = %v, want 0.7", m["temperature"])
		}
		if m["max_tokens"] != 128.0 {
			t.Errorf("max_tokens = %v, want 128", m["max_tokens"])
		}
		if m["reasoning_effort"] != "high" {
			t.Errorf("reasoning_effort = %v, want high", m["reasoning_effort"])
		}
	}, "ok")
	defer srv.Close()

	cfg := Config{BaseURL: srv.URL, APIKey: "sk-test", DefaultModel: "m"}
	_, err := Chat(context.Background(), cfg, "openai", ai.Request{
		Prompt:      "hi",
		System:      "be brief",
		Temperature: 0.7,
		MaxTokens:   128,
		Reasoning:   ai.ReasoningHigh,
	})
	if err != nil {
		t.Fatalf("Chat: %v", err)
	}
}

func TestChatErrors(t *testing.T) {
	if _, err := Chat(context.Background(), Config{}, "x", ai.Request{Prompt: "hi"}); err == nil {
		t.Fatal("expected missing base URL error")
	}
	if _, err := Chat(context.Background(), Config{BaseURL: "http://x"}, "x", ai.Request{Prompt: "hi"}); err == nil {
		t.Fatal("expected missing API key error")
	}
	if _, err := Chat(context.Background(), Config{BaseURL: "http://x", APIKey: "k"}, "x", ai.Request{}); err == nil {
		t.Fatal("expected missing prompt error")
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"bad key"}`))
	}))
	defer srv.Close()
	if _, err := Chat(context.Background(), Config{BaseURL: srv.URL, APIKey: "k", DefaultModel: "m"}, "x", ai.Request{Prompt: "hi"}); err == nil {
		t.Fatal("expected HTTP 401 error")
	}
}
