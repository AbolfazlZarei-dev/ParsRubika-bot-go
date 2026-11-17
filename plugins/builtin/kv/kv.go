// Package kv is a plugin that publishes a small key/value store with
// optional expiry and optional JSON file persistence under the "kv"
// service.
//
// Configuration:
//
//	path            string    JSON file ("" keeps the data in memory only)
//	flush_interval  duration  how often changes are written (default 5s)
package kv

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins/internal/atomicfile"
)

// Store is the published service.
type Store interface {
	Get(key string) (any, bool)
	Set(key string, value any)
	SetTTL(key string, value any, ttl time.Duration)
	Delete(key string)
	Keys(prefix string) []string
}

type item struct {
	Value   any       `json:"v"`
	Expires time.Time `json:"e,omitempty"`
}

// Plugin implements plugins.Plugin and Store.
type Plugin struct {
	plugins.Base
	mu    sync.RWMutex
	data  map[string]item
	path  string
	dirty bool
	now   func() time.Time
}

// New creates the plugin.
func New() *Plugin {
	return &Plugin{Base: plugins.NewBase(plugins.Meta{
		Name: "kv", Version: "1.0.0", Priority: 15,
		Description: "Key/value store service with TTL and persistence",
	}), now: time.Now}
}

// Setup implements plugins.Plugin.
func (p *Plugin) Setup(_ context.Context, h *plugins.Host) error {
	cfg := h.Config()
	p.mu.Lock()
	p.data, p.path = map[string]item{}, cfg.String("path", "")
	p.mu.Unlock()
	if p.path != "" {
		if err := atomicfile.ReadJSON(p.path, &p.data); err != nil {
			return err
		}
		if p.data == nil {
			p.data = map[string]item{}
		}
		h.Every(cfg.Duration("flush_interval", 5*time.Second), func(context.Context) { _ = p.Flush() })
	}
	h.Every(time.Minute, func(context.Context) { p.purge() })
	return h.Provide(plugins.ServiceKV, p)
}

// Teardown implements plugins.Teardowner: a final flush.
func (p *Plugin) Teardown(context.Context) error { return p.Flush() }

// Flush writes pending changes to disk.
func (p *Plugin) Flush() error {
	p.mu.Lock()
	if p.path == "" || !p.dirty {
		p.mu.Unlock()
		return nil
	}
	snap := make(map[string]item, len(p.data))
	for k, v := range p.data {
		snap[k] = v
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

// Get implements Store.
func (p *Plugin) Get(key string) (any, bool) {
	p.mu.RLock()
	it, ok := p.data[key]
	p.mu.RUnlock()
	if !ok || (!it.Expires.IsZero() && !p.now().Before(it.Expires)) {
		return nil, false
	}
	return it.Value, true
}

// Set implements Store.
func (p *Plugin) Set(key string, value any) { p.SetTTL(key, value, 0) }

// SetTTL implements Store; ttl <= 0 means no expiry.
func (p *Plugin) SetTTL(key string, value any, ttl time.Duration) {
	it := item{Value: value}
	if ttl > 0 {
		it.Expires = p.now().Add(ttl)
	}
	p.mu.Lock()
	p.data[key], p.dirty = it, true
	p.mu.Unlock()
}

// Delete implements Store.
func (p *Plugin) Delete(key string) {
	p.mu.Lock()
	delete(p.data, key)
	p.dirty = true
	p.mu.Unlock()
}

// Keys implements Store: live keys with the prefix, sorted.
func (p *Plugin) Keys(prefix string) []string {
	now := p.now()
	p.mu.RLock()
	defer p.mu.RUnlock()
	var out []string
	for k, it := range p.data {
		if len(k) >= len(prefix) && k[:len(prefix)] == prefix && (it.Expires.IsZero() || now.Before(it.Expires)) {
			out = append(out, k)
		}
	}
	sort.Strings(out)
	return out
}

func (p *Plugin) purge() {
	now := p.now()
	p.mu.Lock()
	for k, it := range p.data {
		if !it.Expires.IsZero() && !now.Before(it.Expires) {
			delete(p.data, k)
			p.dirty = true
		}
	}
	p.mu.Unlock()
}
