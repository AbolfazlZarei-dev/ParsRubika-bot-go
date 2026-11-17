package ai

import (
	"context"
	"sync"
)

// Generator is the minimal text-generation callback: prompt in, reply
// out. It is kept as a plain function type so a one-off bot can wire up
// an inline generator (or a mock) without implementing the full
// Provider interface.
//
// Real services should implement Provider instead and be adapted with
// NewGenerator.
type Generator func(ctx context.Context, prompt string) (string, error)

// GeneratorOption configures the Provider→Generator adapter returned by
// NewGenerator.
type GeneratorOption func(*generatorConfig)

type generatorConfig struct {
	system    string
	model     string
	reasoning ReasoningEffort
}

// WithGeneratorSystem prepends a system prompt to every prompt the
// adapter sends to the provider.
func WithGeneratorSystem(system string) GeneratorOption {
	return func(c *generatorConfig) { c.system = system }
}

// WithGeneratorModel overrides the provider's default model for every
// call the adapter makes.
func WithGeneratorModel(model string) GeneratorOption {
	return func(c *generatorConfig) { c.model = model }
}

// WithGeneratorReasoning fixes the reasoning effort for every call the
// adapter makes. Only effective on reasoning-capable providers/models.
func WithGeneratorReasoning(r ReasoningEffort) GeneratorOption {
	return func(c *generatorConfig) { c.reasoning = r }
}

// NewGenerator adapts a Provider into the simple Generator signature.
// The returned function delegates to p.Generate, forwarding ctx, and
// returns Response.Text. It is the bridge between the rich Provider
// interface and the lightweight Generator callback used by Setup.
func NewGenerator(p Provider, opts ...GeneratorOption) Generator {
	if p == nil {
		return func(context.Context, string) (string, error) {
			return "", errNoProvider
		}
	}
	cfg := generatorConfig{}
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	return func(ctx context.Context, prompt string) (string, error) {
		res, err := p.Generate(ctx, Request{
			Model:     cfg.model,
			Prompt:    prompt,
			System:    cfg.system,
			Reasoning: cfg.reasoning,
		})
		if err != nil {
			return "", err
		}
		return res.Text, nil
	}
}

// errNoProvider is returned by a nil-adapted generator.
var errNoProvider = genError("ai: no provider configured")

// genError is a tiny string error so the ai package does not need to
// pull in core/errors (which would create an import cycle through the
// provider sub-packages).
type genError string

func (e genError) Error() string { return string(e) }

// --- global registry ------------------------------------------------------
//
// The registry lets a process register one generator (or provider) at
// startup and have Setup pick it up later without threading it through
// every call site. It is process-global and guarded by a RWMutex.

var (
	registryMu sync.RWMutex
	registry   Generator
)

// RegisterGenerator installs the global generator used by Setup when
// Config.Generator and Config.Provider are both nil. Call this once at
// startup. It overwrites any previously registered generator.
func RegisterGenerator(g Generator) {
	registryMu.Lock()
	defer registryMu.Unlock()
	registry = g
}

// RegisterProvider installs a Provider as the global generator by
// adapting it with NewGenerator. This is the convenience entry point
// for provider packages:
//
//	ai.RegisterProvider(openai.New(openai.WithAPIKey("sk-...")))
func RegisterProvider(p Provider, opts ...GeneratorOption) {
	RegisterGenerator(NewGenerator(p, opts...))
}

// registeredGenerator returns the globally registered generator, or
// nil when none has been set.
func registeredGenerator() Generator {
	registryMu.RLock()
	defer registryMu.RUnlock()
	return registry
}
