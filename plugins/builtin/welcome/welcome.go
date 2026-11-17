// Package welcome is a plugin that greets users when they start the bot
// (the StartedBot event) and answers /start.
//
// The /start handler has plugin priority, so an application that wants
// its own /start simply disables this plugin.
//
// Configuration:
//
//	text     string  the greeting
//	command  string  command name that repeats the greeting (default "start", "" = none)
package welcome

import (
	"context"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot/models"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins"
)

// Plugin implements plugins.Plugin.
type Plugin struct{ plugins.Base }

// New creates the plugin.
func New() *Plugin {
	return &Plugin{plugins.NewBase(plugins.Meta{
		Name: "welcome", Version: "1.0.0",
		Description: "Greeting on start",
	})}
}

// DefaultConfig implements plugins.Defaulter.
func (p *Plugin) DefaultConfig() plugins.Config {
	return plugins.Config{"text": "سلام! 👋 به ربات خوش آمدید.\nبرای دیدن امکانات /help را بفرستید.", "command": "start"}
}

// Setup implements plugins.Plugin.
func (p *Plugin) Setup(_ context.Context, h *plugins.Host) error {
	cfg := h.Config()
	text := cfg.String("text", "")
	h.Use(func(ctx context.Context, b *bot.Bot, u *bot.Update, next func() error) error {
		if u.Type == models.UpdateTypeStartedBot && text != "" {
			if _, err := b.SendMessage(ctx, u.ChatID, text); err != nil {
				return err
			}
		}
		return next()
	})
	if cmd := cfg.String("command", "start"); cmd != "" {
		h.Command(cmd, func(c *bot.Ctx) error { _, err := c.Reply(text); return err })
	}
	return nil
}
