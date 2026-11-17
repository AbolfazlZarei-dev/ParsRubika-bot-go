// Package forcejoin is a plugin that makes users join one or more
// channels before the bot talks to them ("forced subscription").
// Administrators are exempt. If Rubika cannot be asked (API error) the
// user is let through, so an outage never locks everyone out.
//
// Configuration:
//
//	channels  []string  chat IDs of the channels users must have joined
//	links     []string  links shown to users who have not joined
//	message   string    text of the reminder
//	ttl       duration  how long a positive check is remembered (default 5m)
package forcejoin

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins"
)

// Plugin implements plugins.Plugin.
type Plugin struct {
	plugins.Base
	mu     sync.Mutex
	joined map[string]time.Time
	now    func() time.Time
}

// New creates the plugin.
func New() *Plugin {
	return &Plugin{Base: plugins.NewBase(plugins.Meta{
		Name: "forcejoin", Version: "1.0.0", Priority: 40,
		Description: "Forced channel subscription",
	}), now: time.Now}
}

// DefaultConfig implements plugins.Defaulter.
func (p *Plugin) DefaultConfig() plugins.Config {
	return plugins.Config{"ttl": "5m", "message": "برای استفاده از ربات ابتدا در کانال(های) زیر عضو شوید و دوباره تلاش کنید:"}
}

type constErr string

func (e constErr) Error() string { return string(e) }

const errNoChannels = constErr("forcejoin: config \"channels\" must list at least one channel ID")

// Validate implements plugins.Validator.
func (p *Plugin) Validate(c plugins.Config) error {
	if len(c.Strings("channels")) == 0 {
		return errNoChannels
	}
	return nil
}

// Setup implements plugins.Plugin.
func (p *Plugin) Setup(_ context.Context, h *plugins.Host) error {
	cfg := h.Config()
	channels, links := cfg.Strings("channels"), cfg.Strings("links")
	ttl, msg := cfg.Duration("ttl", 5*time.Minute), cfg.String("message", "")
	p.mu.Lock()
	p.joined = map[string]time.Time{}
	p.mu.Unlock()

	h.Use(func(ctx context.Context, b *bot.Bot, u *bot.Update, next func() error) error {
		uid := u.SenderID()
		if uid == "" || !u.IsPrivateChat() || plugins.IsAdmin(h.Manager(), uid) || p.cached(uid, ttl) {
			return next()
		}
		for _, ch := range channels {
			ok, err := b.MemberIsJoined(ctx, ch, uid)
			if err != nil {
				h.Warnf("cannot verify membership of %s in %s: %v (letting the user through)", uid, ch, err)
				return next()
			}
			if !ok {
				text := msg
				if len(links) > 0 {
					text += "\n" + strings.Join(links, "\n")
				}
				_, _ = b.SendMessage(ctx, u.ChatID, text)
				return nil
			}
		}
		p.remember(uid)
		return next()
	})
	h.Every(time.Minute, func(context.Context) { p.sweep(ttl) })
	return nil
}

func (p *Plugin) cached(uid string, ttl time.Duration) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	t, ok := p.joined[uid]
	return ok && p.now().Sub(t) < ttl
}

func (p *Plugin) remember(uid string) {
	p.mu.Lock()
	p.joined[uid] = p.now()
	p.mu.Unlock()
}

func (p *Plugin) sweep(ttl time.Duration) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for k, t := range p.joined {
		if p.now().Sub(t) >= ttl {
			delete(p.joined, k)
		}
	}
}
