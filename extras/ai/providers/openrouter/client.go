// Package openrouter is the OpenRouter provider for the ParsRubika AI
// helper.
//
// OpenRouter is a router over many underlying models, exposed through a
// single OpenAI-compatible endpoint
// (`POST https://openrouter.ai/api/v1/chat/completions`). Because it
// fronts hundreds of models from many vendors, OpenRouter has no single
// "natural" default: the caller must pick a model slug of the form
// `vendor/model` (e.g. `openai/gpt-4o-mini`, `anthropic/claude-opus-5-5`,
// `google/gemini-3.5-flash-lite`). The DefaultModel constant below is
// only a fallback so an empty Request.Model still resolves; always pass
// the model you actually want.
//
// The optional `HTTP-Referer` and `X-Title` headers make your app appear
// on the OpenRouter leaderboards. Set them with WithRankings.
//
//	import "github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/extras/ai/providers/openrouter"
//
//	p := openrouter.New(
//	    openrouter.WithAPIKey("sk-or-v1-..."),
//	    openrouter.WithRankings("https://mybot.example", "My Rubika Bot"),
//	)
//	res, _ := p.Generate(ctx, ai.Request{Model: "anthropic/claude-opus-5-5", Prompt: "Hi"})
package openrouter

import (
	"context"
	"net/http"
	"time"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/extras/ai"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/extras/ai/providers/internal/oai"
)

const (
	// Name is the provider identifier.
	Name = "openrouter"
	// DefaultBaseURL is the OpenRouter API root.
	DefaultBaseURL = "https://openrouter.ai/api/v1"
	// DefaultModel is the fallback model when ai.Request.Model is empty.
	// OpenRouter fronts hundreds of models — always pass the exact
	// `vendor/model` slug you want.
	DefaultModel = "openai/gpt-4o-mini"
)

// Client is the OpenRouter provider. It implements ai.Provider.
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

// WithAPIKey sets the OpenRouter API key (required).
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

// WithModel overrides the default model slug for every Generate call.
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

// WithRankings sets the optional HTTP-Referer (your site URL) and
// X-Title (your app name) headers so your app shows up on OpenRouter's
// leaderboards. Both are optional per OpenRouter's API.
func WithRankings(siteURL, title string) Option {
	return func(c *Client) {
		if c.headers == nil {
			c.headers = map[string]string{}
		}
		if siteURL != "" {
			c.headers["HTTP-Referer"] = siteURL
		}
		if title != "" {
			c.headers["X-Title"] = title
		}
	}
}

// New constructs an OpenRouter Client.
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

// Generate sends a chat-completions request to OpenRouter.
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

// Generator adapts an OpenRouter Client to the simple ai.Generator
// signature.
func Generator(opts ...Option) ai.Generator { return ai.NewGenerator(New(opts...)) }
