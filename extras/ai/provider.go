package ai

import "context"

// ReasoningEffort controls how hard a reasoning-capable model "thinks"
// before answering. Not every provider supports reasoning; providers
// that do (GapGPT, OpenAI o-series, DeepSeek-reasoner, …) map the level
// onto their own wire field and silently ignore it when the chosen
// model cannot reason.
//
// The empty string means "do not send a reasoning field"; use
// ReasoningNone to explicitly ask for no reasoning on providers whose
// default is to reason.
type ReasoningEffort string

const (
	// ReasoningNone disables reasoning entirely.
	ReasoningNone ReasoningEffort = "none"
	// ReasoningLow spends the least effort reasoning.
	ReasoningLow ReasoningEffort = "low"
	// ReasoningMedium is the balanced default for reasoning models.
	ReasoningMedium ReasoningEffort = "medium"
	// ReasoningHigh spends the most effort reasoning.
	ReasoningHigh ReasoningEffort = "high"
)

// Usage reports the token accounting for a single Generate call, when
// the upstream service returns one. Any unknown field stays zero.
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// Request is the input to Provider.Generate. It is provider-agnostic:
// every field is optional except Prompt. Fields a provider does not
// understand are ignored, so the same Request can be replayed across
// providers.
type Request struct {
	// Model is the model id to call. When empty, the provider uses its
	// own DefaultModel. ParsRubika never ships a model catalogue — the
	// caller picks the model, the provider just forwards it.
	Model string

	// Prompt is the user's prompt. Required.
	Prompt string

	// System is an optional system / persona instruction prepended to
	// the conversation. Empty means "no system message".
	System string

	// Temperature is the sampling temperature in [0, 2]. Zero means
	// "unset" — the provider's server-side default applies.
	Temperature float64

	// MaxTokens caps the reply length. Zero means "unset".
	MaxTokens int

	// Reasoning selects the reasoning effort for reasoning-capable
	// models. Empty means "do not send a reasoning field".
	Reasoning ReasoningEffort

	// Extra carries provider-specific knobs that have no home on this
	// struct (e.g. Anthropic `top_p`, Gemini `safetySettings`). It is
	// merged into the outgoing body by providers that consult it.
	Extra map[string]any
}

// Response is the output of Provider.Generate.
type Response struct {
	// Text is the assembled model reply.
	Text string

	// Model is the model id that actually produced the reply (useful
	// when Request.Model was empty and the provider filled in its
	// default).
	Model string

	// Provider is the human-readable provider name
	// (e.g. "openai", "gapgpt").
	Provider string

	// Usage is the token accounting, when the service returned one.
	// Nil means the service did not report usage.
	Usage *Usage

	// Raw is the decoded upstream JSON, for callers that need a field
	// the typed Response does not expose. May be nil.
	Raw map[string]any
}

// Provider is the contract every AI service-provider package implements.
// Each provider (GapGPT, OpenAI, DeepSeek, Anthropic, z.ai, Gemini,
// Groq, OpenRouter, …) lives in its own package under
// `extras/ai/providers/<name>` and registers itself by satisfying this
// interface — ParsRubika never hardcodes the provider list either.
type Provider interface {
	// Name returns the provider's short identifier ("openai",
	// "gapgpt", …). It is surfaced in Response.Provider and in logs.
	Name() string

	// DefaultModel returns the model id used when Request.Model is
	// empty. It is a convenience default chosen by the provider
	// package — callers can always override it per-call.
	DefaultModel() string

	// Generate sends Request to the underlying service and returns the
	// assembled reply. Implementations must honour ctx cancellation.
	Generate(ctx context.Context, req Request) (Response, error)
}
