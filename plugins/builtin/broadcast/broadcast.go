// Package broadcast is a plugin that lets administrators send one
// message to every known chat with /broadcast <text>. It needs the
// access and users plugins.
//
// Configuration:
//
//	command  string    command name (default "broadcast")
//	delay    duration  pause between two sends (default 100ms)
package broadcast

import (
	"context"
	"fmt"
	"time"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins"
)

// Plugin implements plugins.Plugin.
type Plugin struct{ plugins.Base }

// New creates the plugin.
func New() *Plugin {
	return &Plugin{plugins.NewBase(plugins.Meta{
		Name: "broadcast", Version: "1.0.0",
		Description: "Admin broadcast to all chats",
		Requires:    []string{"access", "users"},
	})}
}

// Result summarises a finished broadcast.
type Result struct{ Sent, Failed int }

// Setup implements plugins.Plugin.
func (p *Plugin) Setup(_ context.Context, h *plugins.Host) error {
	cfg := h.Config()
	delay := cfg.Duration("delay", 100*time.Millisecond)
	cmd := cfg.String("command", "broadcast")
	h.Help("/"+cmd+" <متن>", "ارسال پیام همگانی (مدیران)")
	h.Command(cmd, func(c *bot.Ctx) error {
		if !plugins.IsAdmin(h.Manager(), c.SenderID()) {
			return nil
		}
		text := c.ArgsText()
		if text == "" {
			_, err := c.Reply("متن پیام را بعد از دستور بنویسید.")
			return err
		}
		dir, ok := plugins.Service[plugins.UserDirectory](h.Manager(), plugins.ServiceUsers)
		if !ok {
			return fmt.Errorf("broadcast: users service unavailable")
		}
		chats := dir.Chats()
		if _, err := c.Reply(fmt.Sprintf("📣 ارسال به %d گفتگو آغاز شد…", len(chats))); err != nil {
			return err
		}
		chat := c.ChatID()
		// Run in the background so the update queue is not blocked.
		h.Go(func(ctx context.Context) {
			res := Send(ctx, c.Bot, chats, text, delay)
			_, _ = c.Bot.SendMessage(ctx, chat, fmt.Sprintf("✅ ارسال شد: %d — ناموفق: %d", res.Sent, res.Failed))
		})
		return nil
	})
	return nil
}

// Send delivers text to every chat sequentially, pausing delay between
// calls, and stops early when ctx ends.
func Send(ctx context.Context, b *bot.Bot, chats []string, text string, delay time.Duration) Result {
	var r Result
	for _, id := range chats {
		if ctx.Err() != nil {
			break
		}
		if _, err := b.SendMessage(ctx, id, text); err != nil {
			r.Failed++
		} else {
			r.Sent++
		}
		if delay > 0 {
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return r
			}
		}
	}
	return r
}
