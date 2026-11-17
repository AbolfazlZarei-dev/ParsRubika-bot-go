// Package openai is the OpenAI provider for the ParsRubika AI helper.
//
// It speaks the standard OpenAI Chat Completions API
// (`POST https://api.openai.com/v1/chat/completions`) and is reused, via
// the internal `oai` transport, by every other OpenAI-compatible
// provider (DeepSeek, Groq, z.ai, OpenRouter).
//
// ParsRubika does not ship an OpenAI model catalogue. The default model
// below is a convenience; pass any model id you like through
// `ai.Request.Model` (or `openai.WithModel`) the moment OpenAI ships it.
//
//	import "github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/extras/ai"
//	import "github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/extras/ai/providers/openai"
//
//	p := openai.New(openai.WithAPIKey("sk-..."))
//	res, _ := p.Generate(ctx, ai.Request{Prompt: "Hi"})
//	ai.Setup(b, ai.Config{Provider: p})
package openai

import (
	"context"
	"net/http"
	"time"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/extras/ai"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/extras/ai/providers/internal/oai"
)

// Name is the provider identifier surfaced in ai.Response.Provider.
const Name = "openai"

// DefaultBaseURL is the OpenAI API root.
const DefaultBaseURL = "https://api.openai.com/v1"

// DefaultModel is the model used when ai.Request.Model is empty. It is
// a stable, low-cost default; override it per call (or at construction
// with WithModel) for any other OpenAI model.
const DefaultModel = "gpt-4o-mini"

// Client is the OpenAI provider. It implements ai.Provider.
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

// WithAPIKey sets the OpenAI API key (required).
func WithAPIKey(key string) Option {
	return func(c *Client) {
		if key != "" {
			c.apiKey = key
		}
	}
}

// WithBaseURL overrides the default API root (useful for Azure OpenAI or
// self-hosted gateways that mirror the OpenAI shape).
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

// WithHTTPClient injects a pre-configured *http.Client (e.g. one bound
// to a proxy or an HTTP/3 transport).
func WithHTTPClient(h *http.Client) Option {
	return func(c *Client) {
		if h != nil {
			c.httpClient = h
		}
	}
}

// WithHeader adds an extra request header (merged on top of the
// defaults set by oai.Chat).
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

// WithTimeout sets the HTTP timeout used when no custom *http.Client was
// supplied. Zero keeps the default (no timeout — rely on ctx).
func WithTimeout(d time.Duration) Option {
	return func(c *Client) { c.timeout = d }
}

// New constructs an OpenAI Client.
func New(opts ...Option) *Client {
	c := &Client{
		baseURL: DefaultBaseURL,
		model:   DefaultModel,
	}
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

// Generate sends a chat-completions request to OpenAI.
func (c *Client) Generate(ctx context.Context, req ai.Request) (ai.Response, error) {
	return oai.Chat(ctx, oai.Config{
		BaseURL:      c.baseURL,
		APIKey:       c.apiKey,
		DefaultModel: c.model,
		HTTPClient:   c.httpClient,
		Headers:      c.headers,
		UserAgent:    c.userAgent,
	}, Name, req)
}

// Generator is a convenience that builds a Client from opts and adapts
// it to the simple ai.Generator signature. It is the one-liner for
// callers that only want `func(ctx, prompt) (string, error)`.
func Generator(opts ...Option) ai.Generator {
	return ai.NewGenerator(New(opts...))
}
