// Package antispam is a plugin that detects flooding and repeated
// messages and temporarily mutes the offender. Muted users are ignored
// by every later handler. Administrators are exempt.
//
// Configuration:
//
//	max_messages  int       messages allowed inside window (default 6)
//	window        duration  sliding window (default 10s)
//	repeat_limit  int       identical consecutive texts allowed (default 4, 0 = off)
//	mute_for      duration  how long an offender is ignored (default 30s)
//	warn_message  string    sent once when a user gets muted ("" = silent)
package antispam

import (
	"context"
	"sync"
	"time"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins"
)

type userState struct {
	stamps    []time.Time
	lastText  string
	repeats   int
	mutedTill time.Time
	seen      time.Time
}

// Plugin implements plugins.Plugin.
type Plugin struct {
	plugins.Base
	mu     sync.Mutex
	users  map[string]*userState
	max    int
	window time.Duration
	repeat int
	mute   time.Duration
	warn   string
	now    func() time.Time
}

// New creates the plugin.
func New() *Plugin {
	return &Plugin{Base: plugins.NewBase(plugins.Meta{
		Name: "antispam", Version: "1.0.0", Priority: 30,
		Description: "Flood and repeat detection with temporary mutes",
	}), now: time.Now}
}

// DefaultConfig implements plugins.Defaulter.
func (p *Plugin) DefaultConfig() plugins.Config {
	return plugins.Config{
		"max_messages": 6, "window": "10s", "repeat_limit": 4, "mute_for": "30s",
		"warn_message": "⚠️ ارسال پیام‌ها بیش از حد سریع بود؛ برای مدتی نادیده گرفته می‌شوید.",
	}
}

// Setup implements plugins.Plugin.
func (p *Plugin) Setup(_ context.Context, h *plugins.Host) error {
	c := h.Config()
	p.mu.Lock()
	p.max, p.window = c.Int("max_messages", 6), c.Duration("window", 10*time.Second)
	p.repeat, p.mute = c.Int("repeat_limit", 4), c.Duration("mute_for", 30*time.Second)
	p.warn = c.String("warn_message", "")
	p.users = map[string]*userState{}
	p.mu.Unlock()

	h.Use(func(ctx context.Context, b *bot.Bot, u *bot.Update, next func() error) error {
		uid := u.SenderID()
		if uid == "" || u.Message() == nil || plugins.IsAdmin(h.Manager(), uid) {
			return next()
		}
		verdict := p.observe(uid, u.Text())
		switch verdict {
		case allow:
			return next()
		case muteNow:
			if p.warn != "" {
				_, _ = b.SendMessage(ctx, u.ChatID, p.warn)
			}
		}
		return nil
	})
	h.Every(time.Minute, func(context.Context) { p.sweep() })
	return nil
}

type verdict int

const (
	allow verdict = iota
	muteNow
	muted
)

// observe records a message and decides what to do with it.
func (p *Plugin) observe(uid, text string) verdict {
	now := p.now()
	p.mu.Lock()
	defer p.mu.Unlock()
	st := p.users[uid]
	if st == nil {
		st = &userState{}
		p.users[uid] = st
	}
	st.seen = now
	if now.Before(st.mutedTill) {
		return muted
	}
	cut := now.Add(-p.window)
	kept := st.stamps[:0]
	for _, t := range st.stamps {
		if t.After(cut) {
			kept = append(kept, t)
		}
	}
	st.stamps = append(kept, now)

	if text != "" && text == st.lastText {
		st.repeats++
	} else {
		st.repeats, st.lastText = 1, text
	}
	if len(st.stamps) > p.max || (p.repeat > 0 && st.repeats > p.repeat) {
		st.mutedTill = now.Add(p.mute)
		st.stamps, st.repeats = nil, 0
		return muteNow
	}
	return allow
}

// IsMuted reports whether a user is currently muted.
func (p *Plugin) IsMuted(uid string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	st := p.users[uid]
	return st != nil && p.now().Before(st.mutedTill)
}

// Unmute lifts a mute immediately.
func (p *Plugin) Unmute(uid string) {
	p.mu.Lock()
	if st := p.users[uid]; st != nil {
		st.mutedTill = time.Time{}
	}
	p.mu.Unlock()
}

func (p *Plugin) sweep() {
	now := p.now()
	p.mu.Lock()
	defer p.mu.Unlock()
	for k, st := range p.users {
		if now.Sub(st.seen) > 10*time.Minute && now.After(st.mutedTill) {
			delete(p.users, k)
		}
	}
}
