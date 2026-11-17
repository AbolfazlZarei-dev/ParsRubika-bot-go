package plugins

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/core"
)

type entry struct {
	plugin    Plugin
	meta      Meta
	userCfg   Config
	seq       int
	host      *Host
	enabled   bool
	enabledAt time.Time
	lastErr   error
}

// Manager registers plugins and controls their life cycle.
type Manager struct {
	bot *bot.Bot

	lifeMu sync.Mutex // serialises Enable/Disable/Reload

	mu       sync.RWMutex // guards the fields below
	entries  map[string]*entry
	order    []string // enabled plugins, in activation order
	seq      int
	services map[string]any
	preset   map[string]Config
}

// Option customises a Manager.
type Option func(*Manager)

// WithConfigs supplies per-plugin configuration (typically loaded with
// LoadConfigFile). Entries apply to plugins registered later too.
func WithConfigs(c map[string]Config) Option {
	return func(m *Manager) {
		for k, v := range c {
			m.preset[normalize(k)] = v
		}
	}
}

// New creates a Manager bound to b.
func New(b *bot.Bot, opts ...Option) *Manager {
	m := &Manager{
		bot:      b,
		entries:  map[string]*entry{},
		services: map[string]any{},
		preset:   map[string]Config{},
	}
	for _, o := range opts {
		o(m)
	}
	return m
}

// Bot returns the bot the manager is bound to.
func (m *Manager) Bot() *bot.Bot { return m.bot }

func (m *Manager) logf(lv core.Level, format string, a ...any) {
	if m.bot != nil {
		core.Logf(m.bot.Logger(), lv, "plugins: "+format, a...)
	}
}

// ----- registration ------------------------------------------------------

// Register adds a plugin. cfg (optional) overrides the plugin defaults
// and any configuration supplied through WithConfigs.
func (m *Manager) Register(p Plugin, cfg ...Config) error {
	if p == nil {
		return fmt.Errorf("%w: nil plugin", ErrInvalidName)
	}
	meta := p.Meta()
	meta.Name = normalize(meta.Name)
	if !validName(meta.Name) {
		return fmt.Errorf("%w: %q", ErrInvalidName, meta.Name)
	}
	for i, r := range meta.Requires {
		meta.Requires[i] = normalize(r)
	}
	if meta.Priority == 0 {
		meta.Priority = DefaultPriority
	}
	var user Config
	for _, c := range cfg {
		user = Merge(user, c)
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, dup := m.entries[meta.Name]; dup {
		return fmt.Errorf("%w: %s", ErrDuplicate, meta.Name)
	}
	if pre, ok := m.preset[meta.Name]; ok {
		user = Merge(pre, user)
	}
	m.seq++
	m.entries[meta.Name] = &entry{plugin: p, meta: meta, userCfg: user, seq: m.seq}
	return nil
}

// MustRegister is Register that panics on error; convenient in main().
func (m *Manager) MustRegister(p Plugin, cfg ...Config) *Manager {
	if err := m.Register(p, cfg...); err != nil {
		panic(err)
	}
	return m
}

// Unregister disables and removes a plugin.
func (m *Manager) Unregister(ctx context.Context, name string) error {
	name = normalize(name)
	if err := m.Disable(ctx, name); err != nil && !errors.Is(err, ErrNotFound) {
		return err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.entries[name]; !ok {
		return fmt.Errorf("%w: %s", ErrNotFound, name)
	}
	delete(m.entries, name)
	return nil
}

// Configure replaces the user configuration of a plugin. It takes
// effect the next time the plugin is enabled (or with Reload).
func (m *Manager) Configure(name string, cfg Config) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	e, ok := m.entries[normalize(name)]
	if !ok {
		return fmt.Errorf("%w: %s", ErrNotFound, name)
	}
	e.userCfg = cfg.Clone()
	return nil
}

// ----- activation --------------------------------------------------------

// Enable activates the plugins (and, first, everything they require).
func (m *Manager) Enable(ctx context.Context, names ...string) error {
	if m.bot == nil {
		return ErrNoBot
	}
	m.lifeMu.Lock()
	defer m.lifeMu.Unlock()
	var errs []error
	for _, n := range names {
		if err := m.enable(ctx, normalize(n), nil); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// EnableAll activates every registered plugin that is not marked
// Disabled, ordered by priority and then registration order.
func (m *Manager) EnableAll(ctx context.Context) error {
	if m.bot == nil {
		return ErrNoBot
	}
	m.lifeMu.Lock()
	defer m.lifeMu.Unlock()
	m.mu.RLock()
	list := make([]*entry, 0, len(m.entries))
	for _, e := range m.entries {
		if !e.meta.Disabled {
			list = append(list, e)
		}
	}
	m.mu.RUnlock()
	sort.Slice(list, func(i, j int) bool {
		if list[i].meta.Priority != list[j].meta.Priority {
			return list[i].meta.Priority < list[j].meta.Priority
		}
		return list[i].seq < list[j].seq
	})
	var errs []error
	for _, e := range list {
		if err := m.enable(ctx, e.meta.Name, nil); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (m *Manager) lookup(name string) *entry {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.entries[name]
}

// enable is the lock-free core; the caller holds lifeMu.
func (m *Manager) enable(ctx context.Context, name string, stack []string) (err error) {
	e := m.lookup(name)
	if e == nil {
		return fmt.Errorf("%w: %s", ErrNotFound, name)
	}
	if e.enabled {
		return nil
	}
	for _, s := range stack {
		if s == name {
			return fmt.Errorf("%w: cycle %s -> %s", ErrDependency, strings.Join(stack, " -> "), name)
		}
	}
	stack = append(append([]string(nil), stack...), name)
	for _, req := range e.meta.Requires {
		if m.lookup(req) == nil {
			return fmt.Errorf("%w: %q requires unknown plugin %q", ErrDependency, name, req)
		}
		if err := m.enable(ctx, req, stack); err != nil {
			return fmt.Errorf("%w: %q requires %q: %v", ErrDependency, name, req, err)
		}
	}

	cfg := e.userCfg.Clone()
	if d, ok := e.plugin.(Defaulter); ok {
		cfg = Merge(d.DefaultConfig(), e.userCfg)
	}
	if v, ok := e.plugin.(Validator); ok {
		if err := v.Validate(cfg); err != nil {
			e.lastErr = err
			return &SetupError{Plugin: name, Err: fmt.Errorf("invalid config: %w", err)}
		}
	}

	h := newHost(m, name, cfg)
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
		if err != nil {
			_ = h.close(ctx, nil)
			e.lastErr = err
			err = &SetupError{Plugin: name, Err: err}
			m.logf(core.LogError, "cannot enable %s: %v", name, err)
		}
	}()
	if err = e.plugin.Setup(ctx, h); err != nil {
		return err
	}

	m.mu.Lock()
	e.host, e.enabled, e.enabledAt, e.lastErr = h, true, time.Now(), nil
	m.order = append(m.order, name)
	m.mu.Unlock()
	m.logf(core.LogInfo, "enabled %s %s", name, e.meta.Version)
	return nil
}

// Disable deactivates a plugin. Plugins that require it are disabled
// first. Everything the plugin registered is removed from the bot.
func (m *Manager) Disable(ctx context.Context, names ...string) error {
	m.lifeMu.Lock()
	defer m.lifeMu.Unlock()
	var errs []error
	for _, n := range names {
		if err := m.disable(ctx, normalize(n)); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// DisableAll deactivates every plugin, newest first.
func (m *Manager) DisableAll(ctx context.Context) error {
	m.lifeMu.Lock()
	defer m.lifeMu.Unlock()
	m.mu.RLock()
	order := append([]string(nil), m.order...)
	m.mu.RUnlock()
	var errs []error
	for i := len(order) - 1; i >= 0; i-- {
		if err := m.disable(ctx, order[i]); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (m *Manager) disable(ctx context.Context, name string) error {
	e := m.lookup(name)
	if e == nil {
		return fmt.Errorf("%w: %s", ErrNotFound, name)
	}
	if !e.enabled {
		return nil
	}
	// Dependents go first.
	m.mu.RLock()
	var deps []string
	for n, o := range m.entries {
		if o.enabled {
			for _, r := range o.meta.Requires {
				if r == name {
					deps = append(deps, n)
				}
			}
		}
	}
	m.mu.RUnlock()
	var errs []error
	for _, d := range deps {
		if err := m.disable(ctx, d); err != nil {
			errs = append(errs, err)
		}
	}
	var td Teardowner
	if t, ok := e.plugin.(Teardowner); ok {
		td = t
	}
	if err := e.host.close(ctx, td); err != nil {
		errs = append(errs, fmt.Errorf("plugins: teardown %s: %w", name, err))
	}
	m.mu.Lock()
	e.enabled, e.host = false, nil
	for i, n := range m.order {
		if n == name {
			m.order = append(m.order[:i:i], m.order[i+1:]...)
			break
		}
	}
	m.mu.Unlock()
	m.logf(core.LogInfo, "disabled %s", name)
	return errors.Join(errs...)
}

// Reload disables and re-enables a plugin, optionally with a new
// configuration.
func (m *Manager) Reload(ctx context.Context, name string, cfg ...Config) error {
	name = normalize(name)
	if len(cfg) > 0 {
		var c Config
		for _, x := range cfg {
			c = Merge(c, x)
		}
		if err := m.Configure(name, c); err != nil {
			return err
		}
	}
	if err := m.Disable(ctx, name); err != nil {
		return err
	}
	return m.Enable(ctx, name)
}

// Install hooks the manager into the bot life cycle: EnableAll runs when
// the bot starts and DisableAll when it stops.
func (m *Manager) Install() {
	if m.bot == nil {
		return
	}
	m.bot.OnStart(func(ctx context.Context, _ *bot.Bot) error { return m.EnableAll(ctx) })
	m.bot.OnStop(func(ctx context.Context, _ *bot.Bot) error { return m.DisableAll(ctx) })
}

// ----- inspection --------------------------------------------------------

// Info is a read-only snapshot of a plugin.
type Info struct {
	Name        string
	Version     string
	Description string
	Author      string
	Requires    []string
	Priority    int
	Enabled     bool
	EnabledAt   time.Time
	Handlers    int
	Error       string
}

// List returns every registered plugin, ordered by name.
func (m *Manager) List() []Info {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]Info, 0, len(m.entries))
	for _, e := range m.entries {
		in := Info{
			Name: e.meta.Name, Version: e.meta.Version, Description: e.meta.Description,
			Author: e.meta.Author, Requires: append([]string(nil), e.meta.Requires...),
			Priority: e.meta.Priority, Enabled: e.enabled, EnabledAt: e.enabledAt,
		}
		if e.host != nil {
			in.Handlers = e.host.handlerCount()
		}
		if e.lastErr != nil {
			in.Error = e.lastErr.Error()
		}
		out = append(out, in)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// Get returns a registered plugin by name.
func (m *Manager) Get(name string) (Plugin, bool) {
	if e := m.lookup(normalize(name)); e != nil {
		return e.plugin, true
	}
	return nil, false
}

// IsEnabled reports whether a plugin is currently active.
func (m *Manager) IsEnabled(name string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	e := m.entries[normalize(name)]
	return e != nil && e.enabled
}

// Enabled returns the names of active plugins in activation order.
func (m *Manager) Enabled() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return append([]string(nil), m.order...)
}

// HelpEntries returns the documented commands of all active plugins.
func (m *Manager) HelpEntries() []HelpEntry {
	m.mu.RLock()
	var hosts []*Host
	for _, n := range m.order {
		if e := m.entries[n]; e != nil && e.host != nil {
			hosts = append(hosts, e.host)
		}
	}
	m.mu.RUnlock()
	var out []HelpEntry
	for _, h := range hosts {
		out = append(out, h.helpEntries()...)
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].Command < out[j].Command })
	return out
}

// ----- services ----------------------------------------------------------

func (m *Manager) provide(key string, svc any) error {
	if key == "" || svc == nil {
		return errors.New("plugins: service key and value are required")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, dup := m.services[key]; dup {
		return fmt.Errorf("plugins: service %q already provided", key)
	}
	m.services[key] = svc
	return nil
}

func (m *Manager) withdraw(key string) {
	m.mu.Lock()
	delete(m.services, key)
	m.mu.Unlock()
}

// Service returns a published service.
func (m *Manager) Service(key string) (any, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.services[key]
	return s, ok
}

// Service returns the service published under key if it has type T.
func Service[T any](m *Manager, key string) (T, bool) {
	var zero T
	s, ok := m.Service(key)
	if !ok {
		return zero, false
	}
	t, ok := s.(T)
	if !ok {
		return zero, false
	}
	return t, true
}
