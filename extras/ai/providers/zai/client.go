// Package zai is the z.ai (Zhipu GLM) provider for the ParsRubika AI
// helper.
//
// z.ai exposes an OpenAI-compatible chat-completions endpoint at
// `https://api.z.ai/api/paas/v4/chat/completions`, so this package is a
// thin specialisation of the shared `oai` transport. The default model
// is `glm-5.3`; override with any GLM model id z.ai serves
// (`glm-4.6`, `glm-4-flash`, …). ParsRubika does not hardcode a model
// list.
//
// Reasoning models (e.g. those with thinking mode) honour
// ai.Request.Reasoning; pass ai.ReasoningHigh to ask for deep thinking.
//
//	import "github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/extras/ai/providers/zai"
//
//	p := zai.New(zai.WithAPIKey("..."))
//	ai.Setup(b, ai.Config{Provider: p})
package zai

import (
	"context"
	"net/http"
	"time"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/extras/ai"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/extras/ai/providers/internal/oai"
)

const (
	// Name is the provider identifier.
	Name = "zai"
	// DefaultBaseURL is the z.ai OpenAI-compatible API root.
	DefaultBaseURL = "https://api.z.ai/api/paas/v4"
	// DefaultModel is the model used when ai.Request.Model is empty.
	DefaultModel = "glm-5.3"
)

// Client is the z.ai provider. It implements ai.Provider.
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

// WithAPIKey sets the z.ai API key (required).
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

// New constructs a z.ai Client.
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

// Generate sends a chat-completions request to z.ai.
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

// Generator adapts a z.ai Client to the simple ai.Generator signature.
func Generator(opts ...Option) ai.Generator { return ai.NewGenerator(New(opts...)) }
