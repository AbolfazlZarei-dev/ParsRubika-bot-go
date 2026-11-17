// Package requestlog is a plugin that logs every update with its
// processing time and outcome, using the bot's logger (debug level for
// normal updates, warn level for failures and slow handlers).
//
// Configuration:
//
//	slow  duration  updates slower than this are logged as warnings (default 2s)
package requestlog

import (
	"context"
	"time"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins"
)

// Plugin implements plugins.Plugin.
type Plugin struct{ plugins.Base }

// New creates the plugin.
func New() *Plugin {
	return &Plugin{plugins.NewBase(plugins.Meta{
		Name: "requestlog", Version: "1.0.0", Priority: 5,
		Description: "Logs every update with timing",
	})}
}

// Setup implements plugins.Plugin.
func (p *Plugin) Setup(_ context.Context, h *plugins.Host) error {
	slow := h.Config().Duration("slow", 2*time.Second)
	h.Use(func(ctx context.Context, b *bot.Bot, u *bot.Update, next func() error) error {
		start := time.Now()
		err := next()
		d := time.Since(start).Round(time.Microsecond)
		switch {
		case err != nil:
			h.Warnf("%s chat=%s user=%s took %s: %v", u.Type, u.ChatID, u.SenderID(), d, err)
		case d > slow:
			h.Warnf("slow update %s chat=%s took %s", u.Type, u.ChatID, d)
		default:
			h.Debugf("%s chat=%s user=%s took %s", u.Type, u.ChatID, u.SenderID(), d)
		}
		return err
	})
	return nil
}
