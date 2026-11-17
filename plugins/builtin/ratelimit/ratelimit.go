// Package ratelimit is a plugin that limits how fast each user may send
// updates to the bot, using a token bucket per user. Administrators
// (see the access plugin) are exempt.
//
// Configuration:
//
//	rate             float    tokens refilled per second (default 1)
//	burst            int      bucket size (default 5)
//	message          string   reply when a user is throttled ("" = silent)
//	notify_cooldown  duration minimum gap between two replies to one user (default 10s)
package ratelimit

import (
	"context"
	"sync"
	"time"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins"
)

type bucket struct {
	tokens   float64
	last     time.Time
	notified time.Time
}

// Plugin implements plugins.Plugin.
type Plugin struct {
	plugins.Base
	mu      sync.Mutex
	buckets map[string]*bucket
	rate    float64
	burst   float64
	cool    time.Duration
	msg     string
	now     func() time.Time
}

// New creates the plugin.
func New() *Plugin {
	return &Plugin{Base: plugins.NewBase(plugins.Meta{
		Name: "ratelimit", Version: "1.0.0", Priority: 20,
		Description: "Per-user token-bucket rate limiting",
	}), now: time.Now}
}

// DefaultConfig implements plugins.Defaulter.
func (p *Plugin) DefaultConfig() plugins.Config {
	return plugins.Config{"rate": 1.0, "burst": 5, "message": "", "notify_cooldown": "10s"}
}

// Validate implements plugins.Validator.
func (p *Plugin) Validate(c plugins.Config) error {
	if v, ok := c["rate"]; ok && (!isNum(v) || toF(v) <= 0) {
		return errBad("rate must be a positive number")
	}
	if c.Int("burst", 1) < 1 {
		return errBad("burst must be at least 1")
	}
	return nil
}

type errBad string

func (e errBad) Error() string { return string(e) }

func isNum(v any) bool {
	switch v.(type) {
	case int, int64, float64:
		return true
	}
	return false
}

func toF(v any) float64 {
	switch t := v.(type) {
	case int:
		return float64(t)
	case int64:
		return float64(t)
	case float64:
		return t
	}
	return 0
}

// Setup implements plugins.Plugin.
func (p *Plugin) Setup(_ context.Context, h *plugins.Host) error {
	cfg := h.Config()
	p.mu.Lock()
	p.rate = toF(cfg["rate"])
	p.burst = float64(cfg.Int("burst", 5))
	p.cool = cfg.Duration("notify_cooldown", 10*time.Second)
	p.msg = cfg.String("message", "")
	p.buckets = map[string]*bucket{}
	p.mu.Unlock()

	h.Use(func(ctx context.Context, b *bot.Bot, u *bot.Update, next func() error) error {
		uid := u.SenderID()
		if uid == "" || plugins.IsAdmin(h.Manager(), uid) {
			return next()
		}
		ok, notify := p.take(uid)
		if ok {
			return next()
		}
		if notify && p.msg != "" {
			_, _ = b.SendMessage(ctx, u.ChatID, p.msg)
		}
		return nil
	})
	h.Every(time.Minute, func(context.Context) { p.sweep() })
	return nil
}

// take consumes a token for key. notify reports whether the caller may
// be told about the throttling now.
func (p *Plugin) take(key string) (ok, notify bool) {
	now := p.now()
	p.mu.Lock()
	defer p.mu.Unlock()
	b := p.buckets[key]
	if b == nil {
		b = &bucket{tokens: p.burst, last: now}
		p.buckets[key] = b
	}
	b.tokens += now.Sub(b.last).Seconds() * p.rate
	if b.tokens > p.burst {
		b.tokens = p.burst
	}
	b.last = now
	if b.tokens >= 1 {
		b.tokens--
		return true, false
	}
	if now.Sub(b.notified) >= p.cool {
		b.notified = now
		return false, true
	}
	return false, false
}

// sweep drops buckets that have been full (idle) for a while.
func (p *Plugin) sweep() {
	now := p.now()
	p.mu.Lock()
	defer p.mu.Unlock()
	for k, b := range p.buckets {
		if now.Sub(b.last) > 5*time.Minute {
			delete(p.buckets, k)
		}
	}
}
