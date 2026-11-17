# `ai` — Generator

_Source: [`generator.go`](../../../ParsRubika/extras/ai/generator.go)_

## Symbols

### `Generator`

```go
type Generator func(ctx context.Context, prompt string) (string, error)
```

Generator is the minimal text-generation callback: prompt in, reply
out. It is kept as a plain function type so a one-off bot can wire up
an inline generator (or a mock) without implementing the full
Provider interface.

Real services should implement Provider instead and be adapted with
NewGenerator.

### `GeneratorOption`

```go
type GeneratorOption func(*generatorConfig)
```

GeneratorOption configures the Provider→Generator adapter returned by
NewGenerator.

### `WithGeneratorSystem`

```go
func WithGeneratorSystem(system string) GeneratorOption
```

WithGeneratorSystem prepends a system prompt to every prompt the
adapter sends to the provider.

### `WithGeneratorModel`

```go
func WithGeneratorModel(model string) GeneratorOption
```

WithGeneratorModel overrides the provider's default model for every
call the adapter makes.

### `WithGeneratorReasoning`

```go
func WithGeneratorReasoning(r ReasoningEffort) GeneratorOption
```

WithGeneratorReasoning fixes the reasoning effort for every call the
adapter makes. Only effective on reasoning-capable providers/models.

### `NewGenerator`

```go
func NewGenerator(p Provider, opts ...GeneratorOption) Generator
```

NewGenerator adapts a Provider into the simple Generator signature.
The returned function delegates to p.Generate, forwarding ctx, and
returns Response.Text. It is the bridge between the rich Provider
interface and the lightweight Generator callback used by Setup.

### `Error`

```go
func (e genError) Error() string { return string(e) }
```

### `RegisterGenerator`

```go
func RegisterGenerator(g Generator)
```

RegisterGenerator installs the global generator used by Setup when
Config.Generator and Config.Provider are both nil. Call this once at
startup. It overwrites any previously registered generator.

### `RegisterProvider`

```go
func RegisterProvider(p Provider, opts ...GeneratorOption)
```

RegisterProvider installs a Provider as the global generator by
adapting it with NewGenerator. This is the convenience entry point
for provider packages:

ai.RegisterProvider(openai.New(openai.WithAPIKey("sk-...")))
