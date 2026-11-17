// Package oai is the shared transport for OpenAI-compatible chat
// completions providers. It is internal so only the provider packages
// under extras/ai/providers can depend on it.
//
// OpenAI, DeepSeek, Groq, z.ai and OpenRouter all speak the same
// `/chat/completions` dialect: a Bearer token, a `messages` array of
// `{role, content}` objects, and a response whose first choice carries
// the reply in `choices[0].message.content`. Each provider package
// keeps its own base URL, default model and (optional) extra headers,
// then hands the request to oai.Chat to do the JSON round-trip.
//
// The helper never assumes a fixed model catalogue: the model id comes
// from the caller's Request.Model (falling back to the provider's
// DefaultModel), so new models work the moment their vendor ships them.
package oai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/extras/ai"
)

// Config is the per-provider wiring passed to Chat. Every field except
// HTTPClient (which defaults to http.DefaultClient) must be set by the
// provider package.
type Config struct {
	// BaseURL is the API root without a trailing slash, e.g.
	// "https://api.openai.com/v1".
	BaseURL string

	// APIKey is the Bearer token sent in the Authorization header.
	APIKey string

	// DefaultModel is the model id used when ai.Request.Model is empty.
	DefaultModel string

	// HTTPClient is the outbound HTTP client. Nil falls back to
	// http.DefaultClient.
	HTTPClient *http.Client

	// Headers are extra request headers merged on top of the defaults
	// (used by OpenRouter for HTTP-Referer / X-Title ranking, etc.).
	Headers map[string]string

	// UserAgent overrides the default User-Agent header. Empty means
	// "ParsRubika/<version>".
	UserAgent string
}

// message is a single chat-completions message.
type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// chatRequest is the outgoing `/chat/completions` body. Fields are only
// serialised when set, via omitempty, so providers that do not care
// about temperature / max_tokens / reasoning send a minimal payload.
type chatRequest struct {
	Model           string    `json:"model"`
	Messages        []message `json:"messages"`
	Temperature     *float64  `json:"temperature,omitempty"`
	MaxTokens       *int      `json:"max_tokens,omitempty"`
	ReasoningEffort string    `json:"reasoning_effort,omitempty"`
}

// chatResponse is the subset of the OpenAI response shape that oai
// consumes. Unknown fields land in Raw.
type chatResponse struct {
	ID      string `json:"id"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int     `json:"index"`
		Message      message `json:"message"`
		FinishReason string  `json:"finish_reason"`
	} `json:"choices"`
	Usage *ai.Usage `json:"usage,omitempty"`
}

// Chat sends an OpenAI-compatible chat-completions request and returns
// the assembled reply. The provider's Name is stamped onto the
// resulting ai.Response.Provider.
func Chat(ctx context.Context, c Config, name string, req ai.Request) (ai.Response, error) {
	if strings.TrimSpace(c.BaseURL) == "" {
		return ai.Response{}, fmt.Errorf("oai: base URL is required")
	}
	if strings.TrimSpace(c.APIKey) == "" {
		return ai.Response{}, fmt.Errorf("oai: API key is required")
	}
	model := req.Model
	if model == "" {
		model = c.DefaultModel
	}
	if model == "" {
		return ai.Response{}, fmt.Errorf("oai: no model (set Request.Model or Config.DefaultModel)")
	}
	if strings.TrimSpace(req.Prompt) == "" {
		return ai.Response{}, fmt.Errorf("oai: prompt is required")
	}

	body := chatRequest{
		Model: model,
		Messages: []message{
			{Role: "user", Content: req.Prompt},
		},
	}
	if strings.TrimSpace(req.System) != "" {
		// Prepend the system message so the model sees it first.
		body.Messages = append([]message{{Role: "system", Content: req.System}}, body.Messages...)
	}
	if req.Temperature > 0 {
		t := req.Temperature
		body.Temperature = &t
	}
	if req.MaxTokens > 0 {
		m := req.MaxTokens
		body.MaxTokens = &m
	}
	if req.Reasoning != "" {
		body.ReasoningEffort = string(req.Reasoning)
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return ai.Response{}, fmt.Errorf("oai: marshal: %w", err)
	}

	url := strings.TrimRight(c.BaseURL, "/") + "/chat/completions"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return ai.Response{}, fmt.Errorf("oai: new request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)
	ua := c.UserAgent
	if ua == "" {
		ua = "ParsRubika"
	}
	httpReq.Header.Set("User-Agent", ua)
	for k, v := range c.Headers {
		httpReq.Header.Set(k, v)
	}

	client := c.HTTPClient
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Do(httpReq)
	if err != nil {
		return ai.Response{}, fmt.Errorf("oai: request: %w", err)
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(resp.Body, maxBodyBytes))
	if err != nil {
		return ai.Response{}, fmt.Errorf("oai: read body: %w", err)
	}
	if resp.StatusCode >= 400 {
		return ai.Response{}, fmt.Errorf("oai: http %d: %s", resp.StatusCode, truncate(raw, 512))
	}

	var parsed chatResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return ai.Response{}, fmt.Errorf("oai: decode: %w", err)
	}
	if len(parsed.Choices) == 0 {
		return ai.Response{}, fmt.Errorf("oai: response had no choices")
	}
	text := parsed.Choices[0].Message.Content

	var rawMap map[string]any
	_ = json.Unmarshal(raw, &rawMap)

	return ai.Response{
		Text:     text,
		Model:    parsed.Model,
		Provider: name,
		Usage:    parsed.Usage,
		Raw:      rawMap,
	}, nil
}

// maxBodyBytes caps a response body at 8 MiB so a misbehaving endpoint
// cannot exhaust memory.
const maxBodyBytes = 8 << 20

// truncate returns at most n bytes of b as a string, with an ellipsis
// when it was shortened.
func truncate(b []byte, n int) string {
	if len(b) <= n {
		return string(b)
	}
	return string(b[:n]) + "…"
}

// ParseInt is a small helper for providers that read numeric headers.
// Kept here so every provider imports oai for one reason only.
func ParseInt(s string) int { v, _ := strconv.Atoi(strings.TrimSpace(s)); return v }
