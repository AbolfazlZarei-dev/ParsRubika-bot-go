// Package moderation is a plugin with group moderation commands for
// administrators: /ban, /unban, /warn, /warnings and /resetwarns.
// A user who collects max_warns warnings is banned automatically.
//
// Who counts as an administrator: the IDs of the access plugin, or the
// administrators of the group itself.
//
// Configuration:
//
//	max_warns  int  warnings before an automatic ban (default 3)
package moderation

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins"
)

// Plugin implements plugins.Plugin.
type Plugin struct {
	plugins.Base
	mu    sync.Mutex
	warns map[string]int // "chat|user"
}

// New creates the plugin.
func New() *Plugin {
	return &Plugin{Base: plugins.NewBase(plugins.Meta{
		Name: "moderation", Version: "1.0.0",
		Description: "ban / unban / warn commands for admins",
	})}
}

// DefaultConfig implements plugins.Defaulter.
func (p *Plugin) DefaultConfig() plugins.Config { return plugins.Config{"max_warns": 3} }

// Setup implements plugins.Plugin.
func (p *Plugin) Setup(_ context.Context, h *plugins.Host) error {
	maxWarns := h.Config().Int("max_warns", 3)
	p.mu.Lock()
	p.warns = map[string]int{}
	p.mu.Unlock()

	guard := func(fn func(c *bot.Ctx, target string) error) bot.HandlerFunc {
		return func(c *bot.Ctx) error {
			if !p.isAdmin(h, c) {
				return nil
			}
			args := c.Args()
			if len(args) == 0 {
				_, err := c.Reply("شناسه‌ی کاربر را بعد از دستور بنویسید.")
				return err
			}
			return fn(c, args[0])
		}
	}
	h.Help("/ban <id>", "مسدودکردن کاربر در گروه")
	h.Command("ban", guard(func(c *bot.Ctx, t string) error {
		if err := c.Bot.BanChatMember(c, c.ChatID(), t); err != nil {
			return err
		}
		_, err := c.Reply("🚫 کاربر مسدود شد.")
		return err
	}))
	h.Help("/unban <id>", "رفع مسدودی کاربر")
	h.Command("unban", guard(func(c *bot.Ctx, t string) error {
		if err := c.Bot.UnbanChatMember(c, c.ChatID(), t); err != nil {
			return err
		}
		_, err := c.Reply("✅ مسدودی کاربر برداشته شد.")
		return err
	}))
	h.Help("/warn <id> [دلیل]", "اخطار به کاربر")
	h.Command("warn", guard(func(c *bot.Ctx, t string) error {
		n := p.addWarn(c.ChatID(), t)
		if n >= maxWarns {
			p.resetWarns(c.ChatID(), t)
			if err := c.Bot.BanChatMember(c, c.ChatID(), t); err != nil {
				return err
			}
			_, err := c.Reply(fmt.Sprintf("🚫 کاربر به دلیل %d اخطار مسدود شد.", maxWarns))
			return err
		}
		reason := strings.TrimSpace(strings.TrimPrefix(c.ArgsText(), t))
		msg := fmt.Sprintf("⚠️ اخطار %d از %d", n, maxWarns)
		if reason != "" {
			msg += "\nدلیل: " + reason
		}
		_, err := c.Reply(msg)
		return err
	}))
	h.Help("/warnings <id>", "تعداد اخطارهای کاربر")
	h.Command("warnings", guard(func(c *bot.Ctx, t string) error {
		_, err := c.Reply(fmt.Sprintf("تعداد اخطارها: %d از %d", p.warnCount(c.ChatID(), t), maxWarns))
		return err
	}))
	h.Help("/resetwarns <id>", "پاک‌کردن اخطارهای کاربر")
	h.Command("resetwarns", guard(func(c *bot.Ctx, t string) error {
		p.resetWarns(c.ChatID(), t)
		_, err := c.Reply("اخطارها پاک شد.")
		return err
	}))
	return nil
}

func (p *Plugin) isAdmin(h *plugins.Host, c *bot.Ctx) bool {
	if plugins.IsAdmin(h.Manager(), c.SenderID()) {
		return true
	}
	if c.IsPrivate() {
		return false
	}
	ok, err := c.Bot.MemberIsAdmin(c, c.ChatID(), c.SenderID())
	return err == nil && ok
}

func key(chat, user string) string { return chat + "|" + user }

func (p *Plugin) addWarn(chat, user string) int {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.warns[key(chat, user)]++
	return p.warns[key(chat, user)]
}

func (p *Plugin) warnCount(chat, user string) int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.warns[key(chat, user)]
}

func (p *Plugin) resetWarns(chat, user string) {
	p.mu.Lock()
	delete(p.warns, key(chat, user))
	p.mu.Unlock()
}
