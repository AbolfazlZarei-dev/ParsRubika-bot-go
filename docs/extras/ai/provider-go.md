# `ai` — Provider

_Source: [`provider.go`](../../../ParsRubika/extras/ai/provider.go)_

## Symbols

### `ReasoningEffort`

```go
type ReasoningEffort string
```

ReasoningEffort controls how hard a reasoning-capable model "thinks"
before answering. Not every provider supports reasoning; providers
that do (GapGPT, OpenAI o-series, DeepSeek-reasoner, …) map the level
onto their own wire field and silently ignore it when the chosen
model cannot reason.

The empty string means "do not send a reasoning field"; use
ReasoningNone to explicitly ask for no reasoning on providers whose
default is to reason.

### `ReasoningNone`

```go
const ReasoningNone ReasoningEffort = "none"
```

ReasoningNone disables reasoning entirely.

### `ReasoningLow`

```go
const ReasoningLow ReasoningEffort = "low"
```

ReasoningLow spends the least effort reasoning.

### `ReasoningMedium`

```go
const ReasoningMedium ReasoningEffort = "medium"
```

ReasoningMedium is the balanced default for reasoning models.

### `ReasoningHigh`

```go
const ReasoningHigh ReasoningEffort = "high"
```

ReasoningHigh spends the most effort reasoning.

### `Usage`

```go
type Usage struct
```

Usage reports the token accounting for a single Generate call, when
the upstream service returns one. Any unknown field stays zero.

### `Request`

```go
type Request struct
```

Request is the input to Provider.Generate. It is provider-agnostic:
every field is optional except Prompt. Fields a provider does not
understand are ignored, so the same Request can be replayed across
providers.

### `Response`

```go
type Response struct
```

Response is the output of Provider.Generate.

### `Provider`

```go
type Provider interface
```

Provider is the contract every AI service-provider package implements.
Each provider (GapGPT, OpenAI, DeepSeek, Anthropic, z.ai, Gemini,
Groq, OpenRouter, …) lives in its own package under
`extras/ai/providers/<name>` and registers itself by satisfying this
interface — ParsRubika never hardcodes the provider list either.
