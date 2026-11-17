// Package stats is a plugin that exposes the bot's runtime counters
// through an administrator-only /stats command.
//
// Configuration:
//
//	command  string  command name (default "stats")
//	public   bool    let every user run it (default false)
package stats

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins"
)

// Plugin implements plugins.Plugin.
type Plugin struct{ plugins.Base }

// New creates the plugin.
func New() *Plugin {
	return &Plugin{plugins.NewBase(plugins.Meta{
		Name: "stats", Version: "1.0.0",
		Description: "Runtime statistics (/stats)",
	})}
}

// Setup implements plugins.Plugin.
func (p *Plugin) Setup(_ context.Context, h *plugins.Host) error {
	cfg := h.Config()
	public := cfg.Bool("public", false)
	cmd := cfg.String("command", "stats")
	h.Help("/"+cmd, "آمار و وضعیت ربات (مدیران)")
	h.Command(cmd, func(c *bot.Ctx) error {
		if !public && !plugins.IsAdmin(h.Manager(), c.SenderID()) {
			return nil
		}
		_, err := c.Reply(Render(c.Bot.Snapshot(), len(h.Manager().Enabled())))
		return err
	})
	return nil
}

// Render formats a snapshot for humans.
func Render(s bot.StatsSnapshot, plugins int) string {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return fmt.Sprintf("📊 وضعیت ربات\n"+
		"⏱ مدت فعالیت: %s\n"+
		"📥 دریافت‌شده: %d\n"+
		"✅ پردازش‌شده: %d\n"+
		"♻️ تکراری: %d\n"+
		"❌ خطا: %d (panic: %d)\n"+
		"🌐 درخواست API: %d\n"+
		"🧩 پلاگین فعال: %d\n"+
		"🧠 حافظه: %.1f MB — goroutine: %d",
		s.Uptime.Round(time.Second), s.Received, s.Handled, s.Duplicate, s.Failed, s.Panics, s.APICalls,
		plugins, float64(m.Alloc)/(1<<20), runtime.NumGoroutine())
}
