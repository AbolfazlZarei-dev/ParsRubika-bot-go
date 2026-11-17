// Package analytics is a plugin that counts updates by type and
// commands by name and publishes simple counters and gauges (the
// "analytics" service).
package analytics

import (
	"context"
	"sort"
	"sync"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins"
)

// Tracker is the published service.
type Tracker interface {
	Inc(name string, delta float64)
	SetGauge(name string, v float64)
	Counter(name string) float64
	Gauge(name string) float64
	Snapshot() (counters, gauges map[string]float64)
}

// Plugin implements plugins.Plugin and Tracker.
type Plugin struct {
	plugins.Base
	mu       sync.RWMutex
	counters map[string]float64
	gauges   map[string]float64
}

// New creates the plugin.
func New() *Plugin {
	return &Plugin{Base: plugins.NewBase(plugins.Meta{
		Name: "analytics", Version: "1.0.0", Priority: 6,
		Description: "Counters for updates and commands",
	})}
}

// Setup implements plugins.Plugin.
func (p *Plugin) Setup(_ context.Context, h *plugins.Host) error {
	p.mu.Lock()
	p.counters, p.gauges = map[string]float64{}, map[string]float64{}
	p.mu.Unlock()
	h.Use(func(ctx context.Context, b *bot.Bot, u *bot.Update, next func() error) error {
		p.Inc("updates.total", 1)
		p.Inc("updates."+string(u.Type), 1)
		if name := commandOf(u.Text()); name != "" {
			p.Inc("command."+name, 1)
		}
		return next()
	})
	return h.Provide(plugins.ServiceStats, p)
}

func commandOf(text string) string {
	if len(text) < 2 || (text[0] != '/' && text[0] != '!') {
		return ""
	}
	end := len(text)
	for i := 1; i < len(text); i++ {
		if text[i] == ' ' || text[i] == '\n' || text[i] == '@' {
			end = i
			break
		}
	}
	return text[1:end]
}

// Inc implements Tracker.
func (p *Plugin) Inc(name string, d float64) { p.mu.Lock(); p.counters[name] += d; p.mu.Unlock() }

// SetGauge implements Tracker.
func (p *Plugin) SetGauge(name string, v float64) { p.mu.Lock(); p.gauges[name] = v; p.mu.Unlock() }

// Counter implements Tracker.
func (p *Plugin) Counter(name string) float64 {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.counters[name]
}

// Gauge implements Tracker.
func (p *Plugin) Gauge(name string) float64 {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.gauges[name]
}

// Snapshot implements Tracker (copies of both maps).
func (p *Plugin) Snapshot() (map[string]float64, map[string]float64) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	c, g := make(map[string]float64, len(p.counters)), make(map[string]float64, len(p.gauges))
	for k, v := range p.counters {
		c[k] = v
	}
	for k, v := range p.gauges {
		g[k] = v
	}
	return c, g
}

// Names returns the sorted counter names.
func (p *Plugin) Names() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	out := make([]string, 0, len(p.counters))
	for k := range p.counters {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
