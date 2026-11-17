// Package ai integrates ParsRubika bots with any text-generation LLM.
//
// The package is intentionally split across several small files so
// each concern lives in its own unit:
//
//   - doc.go       : this package overview
//   - provider.go  : the Provider interface, Request / Response / Usage
//     types and the ReasoningEffort levels
//   - generator.go : the simple Generator callback, the Provider→Generator
//     adapter and the global generator registry
//   - plugin.go    : the bot-facing Config and Setup that wires an AI
//     reply onto a `/ask`-style command
//   - prompt.go    : the command/prompt extraction helper
//
// ParsRubika does NOT ship a built-in model list. Every supported
// service (GapGPT, OpenAI, DeepSeek, Anthropic, z.ai, Gemini, Groq,
// OpenRouter, …) lives in its own self-contained package under
// `extras/ai/providers/<name>`. Each provider declares its own default
// model and accepts any model string the caller passes in
// `Request.Model`, so new models appear the moment their vendor ships
// them — without a ParsRubika release.
//
// Usage at a glance:
//
//	// 1) Pick a provider (each is its own package).
//	import "github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/extras/ai/providers/openai"
//
//	p := openai.New(openai.WithAPIKey("sk-..."))
//
//	// 2) Either drive it directly …
//	res, err := p.Generate(ctx, ai.Request{Prompt: "Hi"})
//
//	// …or register it with the AI plugin so a Rubika bot answers /ask.
//	ai.Setup(b, ai.Config{Provider: p})
package ai
