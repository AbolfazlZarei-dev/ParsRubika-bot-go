package ai

import (
	"context"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot/models"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot/rules"
)

// Config configures the AI plugin installed by Setup.
//
// A generator can be supplied in three ways, tried in this order:
//  1. Config.Provider — a full Provider implementation (preferred);
//  2. Config.Generator — a plain Generator callback;
//  3. the global registry populated via RegisterProvider/RegisterGenerator.
//
// When Config.Provider is set it is also registered globally so later
// Setup calls (and other components) reuse it without re-passing it.
type Config struct {
	// Provider is the AI service provider. Optional — see the order
	// above.
	Provider Provider

	// Generator is the simple generator callback. Optional.
	Generator Generator

	// Command is the bot command that triggers the AI. Defaults to
	// "ask" (so users type `/ask <prompt>`).
	Command string

	// Prefixes are the command prefixes. Defaults to ["/"].
	Prefixes []string

	// SystemPrompt is prepended to the user's prompt. Useful for
	// giving the model a persona ("You are a friendly Rubika
	// bot…").
	SystemPrompt string

	// Reasoning fixes the reasoning effort passed to the provider for
	// every call. Only effective on reasoning-capable providers/models.
	// Leave empty to let the provider default apply.
	Reasoning ReasoningEffort
}

// Setup registers an AI handler on the bot. The handler fires only for
// messages that are a command matching Config.Command (default
// `/ask <prompt>`). Using a real command rule — instead of a catch-all
// text rule — ensures the AI plugin never shadows other text handlers.
//
// The chosen generator is resolved once at Setup time from
// Config.Provider, Config.Generator or the global registry; if none is
// available the handler replies with a friendly "not configured"
// message instead of panicking.
func Setup(b *bot.Bot, cfg Config) {
	if cfg.Command == "" {
		cfg.Command = "ask"
	}
	if len(cfg.Prefixes) == 0 {
		cfg.Prefixes = []string{"/"}
	}

	// Resolve the generator once. Provider wins, then Generator, then
	// the global registry. When a Provider is given we also adapt it
	// with the configured SystemPrompt/Reasoning and register it
	// globally so subsequent Setup calls reuse it.
	var gen Generator
	switch {
	case cfg.Provider != nil:
		gen = NewGenerator(cfg.Provider,
			WithGeneratorSystem(cfg.SystemPrompt),
			WithGeneratorReasoning(cfg.Reasoning))
		RegisterGenerator(gen)
	case cfg.Generator != nil:
		gen = cfg.Generator
		RegisterGenerator(gen)
	default:
		gen = registeredGenerator()
	}

	b.On([]bot.Rule{rules.NewCommand([]string{cfg.Command}, cfg.Prefixes, true, false)},
		func(ctx context.Context, b *bot.Bot, upd *models.Update) error {
			prompt := extractPrompt(upd.Text(), cfg.Command, cfg.Prefixes)
			if prompt == "" {
				_, _ = b.SendMessage(ctx, upd.ChatID, "Please provide a prompt. Usage: /"+cfg.Command+" <prompt>")
				return nil
			}
			// If the generator came from the registry (gen == nil here
			// but the registry may have been populated later), look it
			// up lazily so RegisterProvider after Setup still works.
			g := gen
			if g == nil {
				g = registeredGenerator()
			}
			if g == nil {
				_, _ = b.SendMessage(ctx, upd.ChatID, "AI generator is not configured.")
				return nil
			}
			fullPrompt := prompt
			if cfg.SystemPrompt != "" && cfg.Provider == nil {
				// Only prepend when the generator is a plain callback
				// that does not honour SystemPrompt itself; the
				// Provider path already folds SystemPrompt in via the
				// adapter.
				fullPrompt = cfg.SystemPrompt + "\n\n" + prompt
			}
			reply, err := g(ctx, fullPrompt)
			if err != nil {
				_, _ = b.SendMessage(ctx, upd.ChatID, "AI error: "+err.Error())
				return err
			}
			_, err = b.SendMessage(ctx, upd.ChatID, reply)
			return err
		})
}
