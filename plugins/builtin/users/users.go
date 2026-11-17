// Package users is a plugin that remembers every chat the bot talks to.
// It marks chats as inactive when the user stops the bot, so a
// broadcast does not waste calls on them. The directory is published as
// the "users" service.
//
// Configuration:
//
//	path            string    JSON file for persistence ("" = memory only)
//	flush_interval  duration  how often changes are written (default 10s)
package users

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot/models"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins/internal/atomicfile"
)

// Record is what is known about one chat.
type Record struct {
	ChatID    string    `json:"chat_id"`
	FirstSeen time.Time `json:"first_seen"`
	LastSeen  time.Time `json:"last_seen"`
	Messages  int       `json:"messages"`
	Stopped   bool      `json:"stopped,omitempty"`
}

// Plugin implements plugins.Plugin and plugins.UserDirectory.
type Plugin struct {
	plugins.Base
	mu    sync.RWMutex
	recs  map[string]*Record
	path  string
	dirty bool
	now   func() time.Time
}

// New creates the plugin.
func New() *Plugin {
	return &Plugin{Base: plugins.NewBase(plugins.Meta{
		Name: "users", Version: "1.0.0", Priority: 25,
		Description: "Remembers chats; powers broadcasts",
	}), now: time.Now}
}

// Setup implements plugins.Plugin.
func (p *Plugin) Setup(_ context.Context, h *plugins.Host) error {
	cfg := h.Config()
	p.mu.Lock()
	p.recs, p.path = map[string]*Record{}, cfg.String("path", "")
	p.mu.Unlock()
	if p.path != "" {
		if err := atomicfile.ReadJSON(p.path, &p.recs); err != nil {
			return err
		}
		if p.recs == nil {
			p.recs = map[string]*Record{}
		}
		h.Every(cfg.Duration("flush_interval", 10*time.Second), func(context.Context) { _ = p.Flush() })
	}
	h.Use(func(ctx context.Context, b *bot.Bot, u *bot.Update, next func() error) error {
		if u.ChatID != "" {
			p.touch(u)
		}
		return next()
	})
	return h.Provide(plugins.ServiceUsers, p)
}

// Teardown implements plugins.Teardowner.
func (p *Plugin) Teardown(context.Context) error { return p.Flush() }

func (p *Plugin) touch(u *bot.Update) {
	now := p.now()
	p.mu.Lock()
	defer p.mu.Unlock()
	r := p.recs[u.ChatID]
	if r == nil {
		r = &Record{ChatID: u.ChatID, FirstSeen: now}
		p.recs[u.ChatID] = r
	}
	r.LastSeen, p.dirty = now, true
	switch u.Type {
	case models.UpdateTypeStoppedBot:
		r.Stopped = true
	case models.UpdateTypeStartedBot:
		r.Stopped = false
	}
	if u.Message() != nil {
		r.Messages++
		r.Stopped = false
	}
}

// Flush writes pending changes to disk.
func (p *Plugin) Flush() error {
	p.mu.Lock()
	if p.path == "" || !p.dirty {
		p.mu.Unlock()
		return nil
	}
	snap := make(map[string]*Record, len(p.recs))
	for k, v := range p.recs {
		c := *v
		snap[k] = &c
	}
	p.dirty = false
	p.mu.Unlock()
	if err := atomicfile.WriteJSON(p.path, snap); err != nil {
		p.mu.Lock()
		p.dirty = true
		p.mu.Unlock()
		return err
	}
	return nil
}

// Chats implements plugins.UserDirectory: active chats, sorted.
func (p *Plugin) Chats() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	out := make([]string, 0, len(p.recs))
	for id, r := range p.recs {
		if !r.Stopped {
			out = append(out, id)
		}
	}
	sort.Strings(out)
	return out
}

// Count implements plugins.UserDirectory.
func (p *Plugin) Count() int { return len(p.Chats()) }

// Get returns the record of a chat.
func (p *Plugin) Get(chatID string) (Record, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	r, ok := p.recs[chatID]
	if !ok {
		return Record{}, false
	}
	return *r, true
}
