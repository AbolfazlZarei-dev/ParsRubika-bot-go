package plugins

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot/rules"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/core"
)

// HelpEntry documents one user-visible command of a plugin.
type HelpEntry struct {
	Plugin      string
	Command     string
	Description string
}

// Host is the plugin's window into the bot. Everything registered
// through it is tracked and removed again when the plugin is disabled,
// so plugins never have to clean up after themselves.
type Host struct {
	m    *Manager
	name string
	bot  *bot.Bot
	cfg  Config

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	mu       sync.Mutex
	handlers []int
	mws      []int
	hooks    []int
	cleanups []func()
	services []string
	help     []HelpEntry
}

func newHost(m *Manager, name string, cfg Config) *Host {
	ctx, cancel := context.WithCancel(context.Background())
	return &Host{m: m, name: name, bot: m.bot, cfg: cfg, ctx: ctx, cancel: cancel}
}

// Name returns the plugin name.
func (h *Host) Name() string { return h.name }

// Bot returns the bot the plugin is attached to.
func (h *Host) Bot() *bot.Bot { return h.bot }

// Manager returns the manager that owns the plugin.
func (h *Host) Manager() *Manager { return h.m }

// Config returns the effective configuration (defaults merged with the
// values supplied at registration time).
func (h *Host) Config() Config { return h.cfg }

// Context is cancelled when the plugin is disabled. Long-running
// goroutines of the plugin should stop when it is done.
func (h *Host) Context() context.Context { return h.ctx }

// ----- logging -----------------------------------------------------------

func (h *Host) logf(lv core.Level, format string, a ...any) {
	if h.bot == nil {
		return
	}
	core.Logf(h.bot.Logger(), lv, "plugin %s: "+format, append([]any{h.name}, a...)...)
}

// Debugf logs at debug level, prefixed with the plugin name.
func (h *Host) Debugf(format string, a ...any) { h.logf(core.LogDebug, format, a...) }

// Infof logs at info level, prefixed with the plugin name.
func (h *Host) Infof(format string, a ...any) { h.logf(core.LogInfo, format, a...) }

// Warnf logs at warning level, prefixed with the plugin name.
func (h *Host) Warnf(format string, a ...any) { h.logf(core.LogWarn, format, a...) }

// Errorf logs at error level, prefixed with the plugin name.
func (h *Host) Errorf(format string, a ...any) { h.logf(core.LogError, format, a...) }

// ----- registration ------------------------------------------------------

func (h *Host) track(list *[]int, id int) int {
	h.mu.Lock()
	*list = append(*list, id)
	h.mu.Unlock()
	return id
}

// On registers a low-level handler. Plugin handlers run before the
// application's ordinary handlers (see bot.PriorityPlugin).
func (h *Host) On(rs []bot.Rule, fn bot.Handler) int {
	return h.track(&h.handlers, h.bot.OnPriority(bot.PriorityPlugin, rs, fn))
}

// Handle registers an ergonomic handler for updates matching all rules.
func (h *Host) Handle(rs []bot.Rule, fn bot.HandlerFunc) int {
	return h.track(&h.handlers, h.bot.HandlePriority(bot.PriorityPlugin, rs, fn))
}

// Command registers a "/name" (or "!name") handler plus aliases.
func (h *Host) Command(name string, fn bot.HandlerFunc, aliases ...string) int {
	return h.Handle([]bot.Rule{bot.CommandRule(name, aliases...)}, fn)
}

// Text registers a handler for exact message texts.
func (h *Host) Text(fn bot.HandlerFunc, exact ...string) int {
	return h.Handle([]bot.Rule{rules.TextIn(exact...)}, fn)
}

// Regex registers a handler for texts matching pattern.
func (h *Host) Regex(pattern string, fn bot.HandlerFunc) int {
	return h.Handle([]bot.Rule{bot.RegexRule(pattern)}, fn)
}

// Button registers a handler for an inline button ID.
func (h *Host) Button(id string, fn bot.HandlerFunc) int {
	return h.Handle([]bot.Rule{rules.NewButton(id, false)}, fn)
}

// Message registers a handler for every text message.
func (h *Host) Message(fn bot.HandlerFunc) int {
	return h.Handle([]bot.Rule{rules.HasText{}}, fn)
}

// Use registers a middleware.
func (h *Host) Use(mw bot.Middleware) int { return h.track(&h.mws, h.bot.Use(mw)) }

// OnStart registers a bot start-up hook.
func (h *Host) OnStart(fn func(ctx context.Context, b *bot.Bot) error) int {
	return h.track(&h.hooks, h.bot.OnStart(fn))
}

// OnStop registers a bot shutdown hook.
func (h *Host) OnStop(fn func(ctx context.Context, b *bot.Bot) error) int {
	return h.track(&h.hooks, h.bot.OnStop(fn))
}

// Help documents a command so the help plugin can list it.
func (h *Host) Help(command, description string) {
	h.mu.Lock()
	h.help = append(h.help, HelpEntry{Plugin: h.name, Command: command, Description: description})
	h.mu.Unlock()
}

// Cleanup schedules fn to run (last in, first out) when the plugin is
// disabled.
func (h *Host) Cleanup(fn func()) {
	h.mu.Lock()
	h.cleanups = append(h.cleanups, fn)
	h.mu.Unlock()
}

// ----- background work ---------------------------------------------------

// Every runs fn repeatedly with the given period until the plugin is
// disabled. A panic in fn is recovered and logged.
func (h *Host) Every(d time.Duration, fn func(ctx context.Context)) {
	if d <= 0 {
		return
	}
	h.wg.Add(1)
	go func() {
		defer h.wg.Done()
		t := time.NewTicker(d)
		defer t.Stop()
		for {
			select {
			case <-h.ctx.Done():
				return
			case <-t.C:
				h.safely(func() { fn(h.ctx) })
			}
		}
	}()
}

// After runs fn once after d unless the plugin is disabled first.
func (h *Host) After(d time.Duration, fn func(ctx context.Context)) {
	h.wg.Add(1)
	go func() {
		defer h.wg.Done()
		t := time.NewTimer(d)
		defer t.Stop()
		select {
		case <-h.ctx.Done():
		case <-t.C:
			h.safely(func() { fn(h.ctx) })
		}
	}()
}

// Go starts fn in a goroutine that is awaited when the plugin is disabled.
func (h *Host) Go(fn func(ctx context.Context)) {
	h.wg.Add(1)
	go func() {
		defer h.wg.Done()
		h.safely(func() { fn(h.ctx) })
	}()
}

func (h *Host) safely(fn func()) {
	defer func() {
		if r := recover(); r != nil {
			h.Errorf("background task panic: %v", r)
		}
	}()
	fn()
}

// ----- services ----------------------------------------------------------

// Provide publishes a service other plugins (and the application) can
// look up with Service. The key must be unique; the service is
// withdrawn when the plugin is disabled.
func (h *Host) Provide(key string, svc any) error {
	if err := h.m.provide(key, svc); err != nil {
		return err
	}
	h.mu.Lock()
	h.services = append(h.services, key)
	h.mu.Unlock()
	return nil
}

// Service looks up a service published by another plugin.
func (h *Host) Service(key string) (any, bool) { return h.m.Service(key) }

// ----- teardown ----------------------------------------------------------

// close undoes everything the plugin registered.
func (h *Host) close(ctx context.Context, td Teardowner) error {
	h.cancel()
	done := make(chan struct{})
	go func() { h.wg.Wait(); close(done) }()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		h.Warnf("background tasks did not stop within 5s")
	}

	var firstErr error
	if td != nil {
		func() {
			defer func() {
				if r := recover(); r != nil {
					firstErr = fmt.Errorf("teardown panic: %v", r)
				}
			}()
			if err := td.Teardown(ctx); err != nil {
				firstErr = err
			}
		}()
	}
	h.mu.Lock()
	cleanups, handlers, mws, hooks, services := h.cleanups, h.handlers, h.mws, h.hooks, h.services
	h.cleanups, h.handlers, h.mws, h.hooks, h.services, h.help = nil, nil, nil, nil, nil, nil
	h.mu.Unlock()

	for i := len(cleanups) - 1; i >= 0; i-- {
		func() {
			defer func() { _ = recover() }()
			cleanups[i]()
		}()
	}
	for _, id := range handlers {
		h.bot.Off(id)
	}
	for _, id := range mws {
		h.bot.RemoveMiddleware(id)
	}
	for _, id := range hooks {
		h.bot.RemoveHook(id)
	}
	for _, key := range services {
		h.m.withdraw(key)
	}
	return firstErr
}

// helpEntries returns a copy of the documented commands.
func (h *Host) helpEntries() []HelpEntry {
	h.mu.Lock()
	defer h.mu.Unlock()
	return append([]HelpEntry(nil), h.help...)
}

// handlerCount returns how many handlers the plugin currently owns.
func (h *Host) handlerCount() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.handlers)
}
