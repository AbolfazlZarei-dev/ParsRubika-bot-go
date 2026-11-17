// Package gemini is the Google Gemini provider for the ParsRubika AI
// helper.
//
// Gemini's API is its own shape: the model id is part of the URL and the
// API key travels as a `?key=` query parameter, not a header. The
// endpoint is
// `POST https://generativelanguage.googleapis.com/v1beta/models/{model}:generateContent`.
// The reply text lives in `candidates[0].content.parts[0].text`.
//
// The default model is `gemini-3.5-flash-lite`, a current general-purpose
// text model. Override with any Gemini model id (`gemini-3.8-flash`,
// `gemini-3.1-pro-preview`, …) — ParsRubika does not hardcode a model
// list. (Note: `gemini-2.5-flash` and `gemini-2.5-pro` are deprecated
// and auto-route away after December 2026; avoid them for new work.)
//
// Thinking-capable Gemini models honour ai.Request.Reasoning via
// `generationConfig.thinkingConfig.thinkingBudget`.
//
//	import "github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/extras/ai/providers/gemini"
//
//	p := gemini.New(gemini.WithAPIKey("..."))
//	ai.Setup(b, ai.Config{Provider: p})
package gemini

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
	Name = "gemini"
	// DefaultBaseURL is the Gemini API root (without the model segment).
	DefaultBaseURL = "https://generativelanguage.googleapis.com/v1beta"
	// DefaultModel is the model used when ai.Request.Model is empty.
	DefaultModel = "gemini-3.5-flash-lite"
)

// Client is the Gemini provider. It implements ai.Provider.
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

// WithAPIKey sets the Gemini API key (sent as the ?key= query param).
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

// WithModel overrides the default Gemini model for every Generate call.
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

// New constructs a Gemini Client.
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

// generateRequest is the outgoing generateContent body.
type generateRequest struct {
	Contents          []geminiContent   `json:"contents"`
	SystemInstruction *geminiContent    `json:"systemInstruction,omitempty"`
	GenerationConfig  *generationConfig `json:"generationConfig,omitempty"`
}

type geminiContent struct {
	Role  string       `json:"role,omitempty"`
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text string `json:"text,omitempty"`
}

type generationConfig struct {
	Temperature     *float64        `json:"temperature,omitempty"`
	MaxOutputTokens *int            `json:"maxOutputTokens,omitempty"`
	ThinkingConfig  *thinkingConfig `json:"thinkingConfig,omitempty"`
}

type thinkingConfig struct {
	ThinkingBudget  int  `json:"thinkingBudget"`
	IncludeThoughts bool `json:"includeThoughts"`
}

// generateResponse is the subset of the Gemini response shape consumed
// here.
type generateResponse struct {
	Candidates []struct {
		Content      geminiContent `json:"content"`
		FinishReason string        `json:"finishReason"`
	} `json:"candidates"`
	UsageMetadata *geminiUsage `json:"usageMetadata,omitempty"`
}

type geminiUsage struct {
	PromptTokenCount     int `json:"promptTokenCount"`
	CandidatesTokenCount int `json:"candidatesTokenCount"`
	TotalTokenCount      int `json:"totalTokenCount"`
}

// Generate sends a generateContent request to Gemini.
func (c *Client) Generate(ctx context.Context, req ai.Request) (ai.Response, error) {
	if strings.TrimSpace(c.apiKey) == "" {
		return ai.Response{}, fmt.Errorf("gemini: API key is required")
	}
	model := req.Model
	if model == "" {
		model = c.model
	}
	if model == "" {
		return ai.Response{}, fmt.Errorf("gemini: no model")
	}
	if strings.TrimSpace(req.Prompt) == "" {
		return ai.Response{}, fmt.Errorf("gemini: prompt is required")
	}

	body := generateRequest{
		Contents: []geminiContent{
			{Role: "user", Parts: []geminiPart{{Text: req.Prompt}}},
		},
	}
	if strings.TrimSpace(req.System) != "" {
		body.SystemInstruction = &geminiContent{Parts: []geminiPart{{Text: req.System}}}
	}

	var gc generationConfig
	set := false
	if req.Temperature > 0 {
		t := req.Temperature
		gc.Temperature = &t
		set = true
	}
	if req.MaxTokens > 0 {
		m := req.MaxTokens
		gc.MaxOutputTokens = &m
		set = true
	}
	// Map reasoning effort onto Gemini's thinkingBudget.
	if req.Reasoning != "" && req.Reasoning != ai.ReasoningNone {
		gc.ThinkingConfig = &thinkingConfig{
			ThinkingBudget:  thinkingBudget(req.Reasoning),
			IncludeThoughts: false,
		}
		set = true
	}
	if set {
		body.GenerationConfig = &gc
	}

	headers := map[string]string{}
	if c.userAgent != "" {
		headers["User-Agent"] = c.userAgent
	}
	for k, v := range c.headers {
		headers[k] = v
	}

	url := strings.TrimRight(c.baseURL, "/") + "/models/" + model + ":generateContent?key=" + c.apiKey
	status, raw, err := httpx.PostJSON(ctx, c.httpClient, url, headers, body)
	if err != nil {
		return ai.Response{}, err
	}
	if status >= 400 {
		return ai.Response{}, fmt.Errorf("gemini: http %d: %s", status, httpx.Truncate(raw, 512))
	}

	var parsed generateResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return ai.Response{}, fmt.Errorf("gemini: decode: %w", err)
	}
	text := ""
	if len(parsed.Candidates) > 0 {
		for _, p := range parsed.Candidates[0].Content.Parts {
			if p.Text != "" {
				text = p.Text
				break
			}
		}
	}

	var usage *ai.Usage
	if parsed.UsageMetadata != nil {
		usage = &ai.Usage{
			PromptTokens:     parsed.UsageMetadata.PromptTokenCount,
			CompletionTokens: parsed.UsageMetadata.CandidatesTokenCount,
			TotalTokens:      parsed.UsageMetadata.TotalTokenCount,
		}
	}

	var rawMap map[string]any
	_ = json.Unmarshal(raw, &rawMap)

	return ai.Response{
		Text:     text,
		Model:    model,
		Provider: Name,
		Usage:    usage,
		Raw:      rawMap,
	}, nil
}

// thinkingBudget maps a ReasoningEffort to a Gemini thinking budget.
func thinkingBudget(r ai.ReasoningEffort) int {
	switch r {
	case ai.ReasoningLow:
		return 2048
	case ai.ReasoningMedium:
		return 8192
	case ai.ReasoningHigh:
		return 24576
	default:
		return 4096
	}
}

// Generator adapts a Gemini Client to the simple ai.Generator signature.
func Generator(opts ...Option) ai.Generator { return ai.NewGenerator(New(opts...)) }
