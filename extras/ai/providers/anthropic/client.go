// Package anthropic is the Anthropic (Claude) provider for the ParsRubika
// AI helper.
//
// Anthropic's Messages API is NOT OpenAI-shaped: it lives at
// `POST https://api.anthropic.com/v1/messages`, authenticates with the
// `x-api-key` header (not a Bearer token) and requires an
// `anthropic-version` header. Requests carry `model`, `max_tokens`, an
// optional top-level `system` string, and a `messages` array. The reply
// text lives in `content[0].text`.
//
// The default model is `claude-opus-5-5` (Anthropic's recommended
// default for most workloads, a fixed id with no date suffix). Override
// with any Claude model id (`claude-sonnet-5-5`, `claude-haiku-5-5`,
// …) — ParsRubika does not hardcode a model list.
//
// Reasoning-capable Claude models honour ai.Request.Reasoning: any level
// other than "none" enables extended thinking with a budget scaled to the
// effort (low=2048, medium=8192, high=32768 tokens).
//
//	import "github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/extras/ai/providers/anthropic"
//
//	p := anthropic.New(anthropic.WithAPIKey("sk-ant-..."))
//	ai.Setup(b, ai.Config{Provider: p})
package anthropic

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
	Name = "anthropic"
	// DefaultBaseURL is the Anthropic API root.
	DefaultBaseURL = "https://api.anthropic.com/v1"
	// DefaultModel is Anthropic's recommended default for most workloads.
	// A fixed id with no date suffix; override per call for any other
	// Claude model.
	DefaultModel = "claude-opus-5-5"
	// APIVersion is the anthropic-version header value. The classic
	// stable version is widely supported across all Claude models.
	APIVersion = "2023-06-01"
	// DefaultMaxTokens is used when ai.Request.MaxTokens is 0, because
	// the Messages API requires a non-zero max_tokens.
	DefaultMaxTokens = 1024
)

// Client is the Anthropic provider. It implements ai.Provider.
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

// WithAPIKey sets the Anthropic API key (sent as x-api-key).
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

// WithModel overrides the default Claude model for every Generate call.
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

// New constructs an Anthropic Client.
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

// messagesRequest is the outgoing Messages API body.
type messagesRequest struct {
	Model     string          `json:"model"`
	MaxTokens int             `json:"max_tokens"`
	System    string          `json:"system,omitempty"`
	Messages  []anthropicMsg  `json:"messages"`
	Thinking  *thinkingConfig `json:"thinking,omitempty"`
}

type anthropicMsg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// thinkingConfig enables extended thinking with a token budget.
type thinkingConfig struct {
	Type         string `json:"type"`
	BudgetTokens int    `json:"budget_tokens"`
}

// messagesResponse is the subset of the Anthropic response shape consumed
// here.
type messagesResponse struct {
	ID      string `json:"id"`
	Model   string `json:"model"`
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	Usage *anthropicUsage `json:"usage,omitempty"`
}

type anthropicUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// Generate sends a Messages API request to Anthropic.
func (c *Client) Generate(ctx context.Context, req ai.Request) (ai.Response, error) {
	if strings.TrimSpace(c.apiKey) == "" {
		return ai.Response{}, fmt.Errorf("anthropic: API key is required")
	}
	model := req.Model
	if model == "" {
		model = c.model
	}
	if model == "" {
		return ai.Response{}, fmt.Errorf("anthropic: no model")
	}
	if strings.TrimSpace(req.Prompt) == "" {
		return ai.Response{}, fmt.Errorf("anthropic: prompt is required")
	}

	maxTokens := req.MaxTokens
	if maxTokens <= 0 {
		maxTokens = DefaultMaxTokens
	}

	body := messagesRequest{
		Model:     model,
		MaxTokens: maxTokens,
		System:    req.System,
		Messages:  []anthropicMsg{{Role: "user", Content: req.Prompt}},
	}

	// Map reasoning effort onto Anthropic extended thinking. "none" or
	// empty disables it; any other level enables it with a scaled budget.
	if req.Reasoning != "" && req.Reasoning != ai.ReasoningNone {
		budget := thinkingBudget(req.Reasoning)
		// max_tokens must exceed the thinking budget.
		if maxTokens <= budget {
			maxTokens = budget + 1024
		}
		body.MaxTokens = maxTokens
		body.Thinking = &thinkingConfig{Type: "enabled", BudgetTokens: budget}
	}

	headers := map[string]string{
		"x-api-key":         c.apiKey,
		"anthropic-version": APIVersion,
	}
	if c.userAgent != "" {
		headers["User-Agent"] = c.userAgent
	}
	for k, v := range c.headers {
		headers[k] = v
	}

	url := strings.TrimRight(c.baseURL, "/") + "/messages"
	status, raw, err := httpx.PostJSON(ctx, c.httpClient, url, headers, body)
	if err != nil {
		return ai.Response{}, err
	}
	if status >= 400 {
		return ai.Response{}, fmt.Errorf("anthropic: http %d: %s", status, httpx.Truncate(raw, 512))
	}

	var parsed messagesResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return ai.Response{}, fmt.Errorf("anthropic: decode: %w", err)
	}
	text := ""
	for _, blk := range parsed.Content {
		if blk.Type == "text" && blk.Text != "" {
			text = blk.Text
			break
		}
	}

	var usage *ai.Usage
	if parsed.Usage != nil {
		usage = &ai.Usage{
			PromptTokens:     parsed.Usage.InputTokens,
			CompletionTokens: parsed.Usage.OutputTokens,
			TotalTokens:      parsed.Usage.InputTokens + parsed.Usage.OutputTokens,
		}
	}

	var rawMap map[string]any
	_ = json.Unmarshal(raw, &rawMap)

	return ai.Response{
		Text:     text,
		Model:    parsed.Model,
		Provider: Name,
		Usage:    usage,
		Raw:      rawMap,
	}, nil
}

// thinkingBudget maps a ReasoningEffort to an extended-thinking token
// budget.
func thinkingBudget(r ai.ReasoningEffort) int {
	switch r {
	case ai.ReasoningLow:
		return 2048
	case ai.ReasoningMedium:
		return 8192
	case ai.ReasoningHigh:
		return 32768
	default:
		return 4096
	}
}

// Generator adapts an Anthropic Client to the simple ai.Generator
// signature.
func Generator(opts ...Option) ai.Generator { return ai.NewGenerator(New(opts...)) }
