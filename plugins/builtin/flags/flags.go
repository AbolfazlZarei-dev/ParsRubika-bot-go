// Package flags is a plugin that publishes feature flags (the "flags"
// service): global switches plus per-user overrides, handy for staged
// roll-outs and beta features.
//
// Configuration:
//
//	flags  map[string]bool  initial flag values
package flags

import (
	"context"
	"sync"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins"
)

// Flags is the published service.
type Flags interface {
	IsEnabled(name string) bool
	IsEnabledFor(name, userID string) bool
	Set(name string, on bool)
	SetFor(name, userID string, on bool)
}

// Plugin implements plugins.Plugin and Flags.
type Plugin struct {
	plugins.Base
	mu     sync.RWMutex
	global map[string]bool
	user   map[string]map[string]bool
}

// New creates the plugin.
func New() *Plugin {
	return &Plugin{Base: plugins.NewBase(plugins.Meta{
		Name: "flags", Version: "1.0.0", Priority: 15,
		Description: "Feature flags with per-user overrides",
	})}
}

// Setup implements plugins.Plugin.
func (p *Plugin) Setup(_ context.Context, h *plugins.Host) error {
	p.mu.Lock()
	p.global, p.user = map[string]bool{}, map[string]map[string]bool{}
	switch raw := h.Config()["flags"].(type) {
	case map[string]any:
		for k, v := range raw {
			if b, ok := v.(bool); ok {
				p.global[k] = b
			}
		}
	case map[string]bool:
		for k, v := range raw {
			p.global[k] = v
		}
	}
	p.mu.Unlock()
	return h.Provide(plugins.ServiceFlags, p)
}

// IsEnabled implements Flags.
func (p *Plugin) IsEnabled(name string) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.global[name]
}

// IsEnabledFor implements Flags: a per-user override wins over the
// global value.
func (p *Plugin) IsEnabledFor(name, userID string) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if v, ok := p.user[name][userID]; ok {
		return v
	}
	return p.global[name]
}

// Set implements Flags.
func (p *Plugin) Set(name string, on bool) { p.mu.Lock(); p.global[name] = on; p.mu.Unlock() }

// SetFor implements Flags.
func (p *Plugin) SetFor(name, userID string, on bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.user[name] == nil {
		p.user[name] = map[string]bool{}
	}
	p.user[name][userID] = on
}
