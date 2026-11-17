// Package gapgpt is the GapGPT provider for the ParsRubika AI helper.
//
// GapGPT (https://api.gapgpt.app) is an Iranian-hosted gateway that
// fronts OpenAI, Gemini, Claude, DeepSeek and many other models behind
// a single OpenAI Responses-style endpoint:
//
//	POST https://api.gapgpt.app/v1/responses
//	Authorization: Bearer YOUR_API_KEY
//	{"model": "...", "input": "...", "reasoning": {"effort": "none|low|medium|high"}}
//
// The default model is `gapgpt-deepseek-v4.1-flash` (GapGPT's own
// DeepSeek hosting). Override with any model id GapGPT serves —
// ParsRubika does not hardcode a model list.
//
// Reasoning effort maps directly onto GapGPT's `reasoning.effort` field:
// pass ai.ReasoningNone to disable reasoning, or low/medium/high to dial
// it up. (GapGPT notes that for some models — e.g. Gemma — reasoning is
// binary: low, medium and high all turn it on.)
//
//	import "github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/extras/ai/providers/gapgpt"
//
//	p := gapgpt.New(gapgpt.WithAPIKey("..."))
//	res, _ := p.Generate(ctx, ai.Request{
//	    Prompt: "Write a short bedtime story",
//	    Reasoning: ai.ReasoningLow,
//	})
//	ai.Setup(b, ai.Config{Provider: p})
package gapgpt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/extras/ai"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/extras/ai/providers/internal/httpx"
)

const (
	// Name is the provider identifier.
	Name = "gapgpt"
	// DefaultBaseURL is the GapGPT API root.
	DefaultBaseURL = "https://api.gapgpt.app/v1"
	// DefaultModel is GapGPT's hosted DeepSeek flash model.
	DefaultModel = "gapgpt-deepseek-v4.1-flash"
)

// Client is the GapGPT provider. It implements ai.Provider.
type Client struct {
	apiKey     string
	baseURL    string
	model      string
	httpClient *http.Client
	headers    map[string]string
	userAgent  string
	timeout    time.Duration
}

// Option configures a Client at construction time.
type Option func(*Client)

// WithAPIKey sets the GapGPT API key (Bearer token).
func WithAPIKey(key string) Option {
	return func(c *Client) {
		if key != "" {
			c.apiKey = key
		}
	}
}

// WithBaseURL overrides the default API root.
func WithBaseURL(url string) Option {
	return func(c *Client) {
		if url != "" {
			c.baseURL = url
		}
	}
}

// WithModel overrides the default model for every Generate call.
func WithModel(model string) Option {
	return func(c *Client) {
		if model != "" {
			c.model = model
		}
	}
}

// WithHTTPClient injects a pre-configured *http.Client.
func WithHTTPClient(h *http.Client) Option {
	return func(c *Client) {
		if h != nil {
			c.httpClient = h
		}
	}
}

// WithHeader adds an extra request header.
func WithHeader(key, value string) Option {
	return func(c *Client) {
		if c.headers == nil {
			c.headers = map[string]string{}
		}
		c.headers[key] = value
	}
}

// WithUserAgent overrides the default User-Agent header.
func WithUserAgent(ua string) Option {
	return func(c *Client) {
		if ua != "" {
			c.userAgent = ua
		}
	}
}

// WithTimeout sets the HTTP timeout when no custom *http.Client is set.
func WithTimeout(d time.Duration) Option { return func(c *Client) { c.timeout = d } }

// New constructs a GapGPT Client.
func New(opts ...Option) *Client {
	c := &Client{baseURL: DefaultBaseURL, model: DefaultModel}
	for _, opt := range opts {
		if opt != nil {
			opt(c)
		}
	}
	if c.httpClient == nil && c.timeout > 0 {
		c.httpClient = &http.Client{Timeout: c.timeout}
	}
	return c
}

// Name returns the provider identifier.
func (c *Client) Name() string { return Name }

// DefaultModel returns the model used when ai.Request.Model is empty.
func (c *Client) DefaultModel() string { return c.model }

// responsesRequest is the outgoing GapGPT Responses body.
type responsesRequest struct {
	Model     string          `json:"model"`
	Input     string          `json:"input"`
	Reasoning *reasoningBlock `json:"reasoning,omitempty"`
}

// reasoningBlock carries the effort level. GapGPT accepts
// "none" | "low" | "medium" | "high".
type reasoningBlock struct {
	Effort string `json:"effort"`
}

// Generate sends a Responses request to GapGPT.
func (c *Client) Generate(ctx context.Context, req ai.Request) (ai.Response, error) {
	if strings.TrimSpace(c.apiKey) == "" {
		return ai.Response{}, fmt.Errorf("gapgpt: API key is required")
	}
	model := req.Model
	if model == "" {
		model = c.model
	}
	if model == "" {
		return ai.Response{}, fmt.Errorf("gapgpt: no model")
	}
	if strings.TrimSpace(req.Prompt) == "" {
		return ai.Response{}, fmt.Errorf("gapgpt: prompt is required")
	}

	input := req.Prompt
	if strings.TrimSpace(req.System) != "" {
		// GapGPT's Responses API takes a single `input` string; fold the
		// system prompt in front of the user prompt.
		input = req.System + "\n\n" + req.Prompt
	}

	body := responsesRequest{Model: model, Input: input}
	// GapGPT distinguishes "none" (no reasoning) from the empty case
	// (provider default). Forward the effort whenever the caller set it.
	if req.Reasoning != "" {
		body.Reasoning = &reasoningBlock{Effort: string(req.Reasoning)}
	}

	headers := map[string]string{
		"Authorization": "Bearer " + c.apiKey,
	}
	if c.userAgent != "" {
		headers["User-Agent"] = c.userAgent
	}
	for k, v := range c.headers {
		headers[k] = v
	}

	url := strings.TrimRight(c.baseURL, "/") + "/responses"
	status, raw, err := httpx.PostJSON(ctx, c.httpClient, url, headers, body)
	if err != nil {
		return ai.Response{}, err
	}
	if status >= 400 {
		return ai.Response{}, fmt.Errorf("gapgpt: http %d: %s", status, httpx.Truncate(raw, 512))
	}

	var parsed map[string]any
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return ai.Response{}, fmt.Errorf("gapgpt: decode: %w", err)
	}

	text := extractResponsesText(parsed)
	modelOut, _ := parsed["model"].(string)
	if modelOut == "" {
		modelOut = model
	}

	usage := extractUsage(parsed)

	return ai.Response{
		Text:     text,
		Model:    modelOut,
		Provider: Name,
		Usage:    usage,
		Raw:      parsed,
	}, nil
}

// extractResponsesText pulls the reply text out of a Responses-style
// payload. It tries, in order:
//  1. output_text (top-level convenience string)
//  2. output[].content[].text (the structured message shape)
//  3. choices[0].message.content (chat-completions fallback some
//     gateways return instead)
func extractResponsesText(root map[string]any) string {
	if s, ok := root["output_text"].(string); ok && s != "" {
		return s
	}
	if outs, ok := root["output"].([]any); ok {
		for _, o := range outs {
			om, _ := o.(map[string]any)
			if contents, ok := om["content"].([]any); ok {
				for _, c := range contents {
					cm, _ := c.(map[string]any)
					if t, _ := cm["type"].(string); t == "output_text" || t == "text" {
						if s, _ := cm["text"].(string); s != "" {
							return s
						}
					}
				}
			}
			// Some shapes put text directly on the output item.
			if s, _ := om["text"].(string); s != "" {
				return s
			}
		}
	}
	if choices, ok := root["choices"].([]any); ok && len(choices) > 0 {
		cm, _ := choices[0].(map[string]any)
		msg, _ := cm["message"].(map[string]any)
		if s, _ := msg["content"].(string); s != "" {
			return s
		}
	}
	return ""
}

// extractUsage reads a usage block whether the keys are named in the
// OpenAI Responses style (input_tokens/output_tokens) or the
// chat-completions style (prompt_tokens/completion_tokens).
func extractUsage(root map[string]any) *ai.Usage {
	u, ok := root["usage"].(map[string]any)
	if !ok {
		return nil
	}
	get := func(keys ...string) int {
		for _, k := range keys {
			if v, ok := u[k]; ok {
				switch n := v.(type) {
				case float64:
					return int(n)
				case int:
					return n
				}
			}
		}
		return 0
	}
	prompt := get("input_tokens", "prompt_tokens")
	completion := get("output_tokens", "completion_tokens")
	total := get("total_tokens")
	if total == 0 {
		total = prompt + completion
	}
	if prompt == 0 && completion == 0 {
		return nil
	}
	return &ai.Usage{PromptTokens: prompt, CompletionTokens: completion, TotalTokens: total}
}

// Generator adapts a GapGPT Client to the simple ai.Generator signature.
func Generator(opts ...Option) ai.Generator { return ai.NewGenerator(New(opts...)) }
