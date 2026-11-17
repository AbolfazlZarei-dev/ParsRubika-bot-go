package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/core"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/text"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/transport"
)

// Handler is the callback signature used by every registration helper.
type Handler = func(ctx context.Context, b *Bot, upd *Update) error

// Middleware wraps the dispatch of an update. Calling next runs the
// remaining middlewares and, finally, the matching handler; not calling
// it stops the update from reaching any handler.
type Middleware func(ctx context.Context, b *Bot, upd *Update, next func() error) error

// handlerEntry pairs a slice of rules with a callback.
type handlerEntry struct {
	id       int
	prio     int
	rules    []Rule
	callback Handler
}

type middlewareEntry struct {
	id int
	fn Middleware
}

type hookEntry struct {
	id int
	fn func(ctx context.Context, b *Bot) error
}

// Bot is the bot API client. A Bot is safe for concurrent use; handlers,
// middlewares and hooks may be registered or removed while it runs.
type Bot struct {
	token   string
	opts    Options
	baseURL string

	transport *transport.Transport
	formatter text.Formatter
	logger    Logger

	mu            sync.RWMutex
	handlers      []handlerEntry
	middlewares   []middlewareEntry
	startHandlers []hookEntry
	stopHandlers  []hookEntry
	nextID        int

	// Polling bookkeeping.
	offsetMu     sync.Mutex
	nextOffsetID string
	offsetFile   string
	hadOffset    bool
	skipPending  atomic.Bool
	seen         *seenSet
	rateLimiter  *transport.RateLimiter

	disp *dispatcher

	// Lifecycle.
	running   atomic.Bool
	lifeMu    sync.Mutex
	lifeCtx   context.Context
	lifeStop  context.CancelFunc
	startedAt time.Time

	stats Stats
	me    meCache
}

// Stats holds lightweight runtime counters.
type Stats struct {
	Received  atomic.Uint64
	Handled   atomic.Uint64
	Duplicate atomic.Uint64
	Failed    atomic.Uint64
	Panics    atomic.Uint64
	APICalls  atomic.Uint64
}

// StatsSnapshot is a point-in-time copy of Stats.
type StatsSnapshot struct {
	Received, Handled, Duplicate, Failed, Panics, APICalls uint64
	Uptime                                                 time.Duration
}

// New constructs a Bot using the supplied token. Options may be given
// either as an Options value or through the With* helpers.
func New(token string, opts ...Option) *Bot {
	o := DefaultOptions()
	for _, opt := range opts {
		if opt != nil {
			opt.apply(&o)
		}
	}
	if o.BaseURL == "" {
		o.BaseURL = BaseURL
	}
	if o.Concurrency < 1 {
		o.Concurrency = 1
	}
	url := fmt.Sprintf(o.BaseURL, token)
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}
	var lg Logger = o.Logger
	if lg == nil {
		lg = stdLogger{log.New(os.Stderr, "[parsrubika-bot] ", log.LstdFlags|log.Lmsgprefix)}
	}
	b := &Bot{
		token:       token,
		opts:        o,
		baseURL:     url,
		seen:        newSeenSet(o.DedupSize),
		rateLimiter: transport.NewRateLimiter(o.RateLimit),
		logger:      lg,
		disp:        newDispatcher(o.Concurrency),
	}
	t := transport.New(
		transport.WithUserAgent(o.UserAgent),
		transport.WithTimeout(o.Timeout),
		transport.WithMaxRetries(o.MaxRetries),
	)
	if o.HTTPClient != nil {
		t.SetHTTPClient(o.HTTPClient)
	}
	b.transport = t
	if o.PersistOffset {
		b.offsetFile = filepath.Join(o.OffsetDir, ".parsrubika_offset_"+shortHash(token))
		if data, err := os.ReadFile(b.offsetFile); err == nil {
			if v := strings.TrimSpace(string(data)); v != "" {
				b.nextOffsetID = v
				b.hadOffset = true
			}
		}
	}
	b.skipPending.Store(o.SkipBacklog && !b.hadOffset)
	return b
}

// SetHTTPClient lets callers swap in a custom HTTP client.
func (b *Bot) SetHTTPClient(c *http.Client) {
	if c != nil {
		b.transport.SetHTTPClient(c)
	}
}

// Token returns the bot token.
func (b *Bot) Token() string { return b.token }

// BaseURL returns the resolved bot API base URL.
func (b *Bot) BaseURL() string { return b.baseURL }

// Logger returns the logger used by the bot.
func (b *Bot) Logger() Logger { return b.logger }

// Options returns a copy of the effective configuration.
func (b *Bot) Options() Options { return b.opts }

// Snapshot returns the current runtime counters.
func (b *Bot) Snapshot() StatsSnapshot {
	s := StatsSnapshot{
		Received:  b.stats.Received.Load(),
		Handled:   b.stats.Handled.Load(),
		Duplicate: b.stats.Duplicate.Load(),
		Failed:    b.stats.Failed.Load(),
		Panics:    b.stats.Panics.Load(),
		APICalls:  b.stats.APICalls.Load(),
	}
	if !b.startedAt.IsZero() {
		s.Uptime = time.Since(b.startedAt)
	}
	return s
}

func (b *Bot) newID() int {
	b.nextID++
	return b.nextID
}

// ===== Registration =====================================================

// On registers a handler that fires when ALL of the supplied rules
// match. The returned ID can be passed to Off.
func (b *Bot) On(rules []Rule, handler Handler) int { return b.OnPriority(0, rules, handler) }

// Handler priorities. Handlers with a higher priority are tried first;
// handlers of equal priority run in registration order.
const (
	// PriorityPlugin is used for handlers registered by plugins, so a
	// catch-all handler of the application cannot swallow their commands.
	PriorityPlugin = 50
	// PriorityFallback is used by Fallback: it always runs last.
	PriorityFallback = -100
)

// OnPriority is On with an explicit priority (default 0).
func (b *Bot) OnPriority(prio int, rules []Rule, handler Handler) int {
	b.mu.Lock()
	defer b.mu.Unlock()
	id := b.newID()
	e := handlerEntry{id: id, prio: prio, rules: rules, callback: handler}
	i := len(b.handlers)
	for j, h := range b.handlers {
		if h.prio < prio {
			i = j
			break
		}
	}
	b.handlers = append(b.handlers, handlerEntry{})
	copy(b.handlers[i+1:], b.handlers[i:])
	b.handlers[i] = e
	return id
}

// Off removes a handler registered with On. It reports whether the
// handler existed.
func (b *Bot) Off(id int) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	for i, h := range b.handlers {
		if h.id == id {
			b.handlers = append(b.handlers[:i:i], b.handlers[i+1:]...)
			return true
		}
	}
	return false
}

// OnStart registers a startup hook and returns its ID.
func (b *Bot) OnStart(h func(ctx context.Context, b *Bot) error) int {
	b.mu.Lock()
	defer b.mu.Unlock()
	id := b.newID()
	b.startHandlers = append(b.startHandlers, hookEntry{id: id, fn: h})
	return id
}

// OnStop registers a shutdown hook and returns its ID.
func (b *Bot) OnStop(h func(ctx context.Context, b *Bot) error) int {
	b.mu.Lock()
	defer b.mu.Unlock()
	id := b.newID()
	b.stopHandlers = append(b.stopHandlers, hookEntry{id: id, fn: h})
	return id
}

// RemoveHook removes a hook registered with OnStart or OnStop.
func (b *Bot) RemoveHook(id int) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	rm := func(list []hookEntry) ([]hookEntry, bool) {
		for i, h := range list {
			if h.id == id {
				return append(list[:i:i], list[i+1:]...), true
			}
		}
		return list, false
	}
	var ok bool
	if b.startHandlers, ok = rm(b.startHandlers); ok {
		return true
	}
	b.stopHandlers, ok = rm(b.stopHandlers)
	return ok
}

// Use registers a middleware and returns its ID. Middlewares run in
// registration order, outermost first.
func (b *Bot) Use(m Middleware) int {
	b.mu.Lock()
	defer b.mu.Unlock()
	id := b.newID()
	b.middlewares = append(b.middlewares, middlewareEntry{id: id, fn: m})
	return id
}

// RemoveMiddleware removes a middleware registered with Use.
func (b *Bot) RemoveMiddleware(id int) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	for i, m := range b.middlewares {
		if m.id == id {
			b.middlewares = append(b.middlewares[:i:i], b.middlewares[i+1:]...)
			return true
		}
	}
	return false
}

// ===== Lifecycle ========================================================

// Start runs the bot in polling mode. It blocks until ctx is
// cancelled, Stop is called, or a fatal error occurs.
func (b *Bot) Start(ctx context.Context) error { return b.Run(ctx, "") }

// Run starts the bot in polling mode, or in webhook mode when
// webhookURL is non-empty (the optional argument is "host:port" of the
// local listener).
func (b *Bot) Run(ctx context.Context, webhookURL string, hostAndPort ...string) error {
	if !b.running.CompareAndSwap(false, true) {
		return fmt.Errorf("bot: already running")
	}
	defer b.running.Store(false)

	b.lifeMu.Lock()
	lctx, cancel := context.WithCancel(ctx)
	b.lifeCtx, b.lifeStop = lctx, cancel
	b.startedAt = time.Now()
	b.lifeMu.Unlock()
	defer cancel()

	if webhookURL != "" {
		host, port := "0.0.0.0", 8080
		if len(hostAndPort) > 0 && hostAndPort[0] != "" {
			if h, p, err := net.SplitHostPort(hostAndPort[0]); err == nil {
				if h != "" {
					host = h
				}
				port = parseIntDefault(p, 8080)
			} else {
				host = hostAndPort[0]
			}
		}
		return b.runWebhook(lctx, webhookURL, host, port)
	}
	return b.runPolling(lctx)
}

// IsRunning reports whether the bot is currently running.
func (b *Bot) IsRunning() bool { return b.running.Load() }

// Stop asks the polling loop or webhook server to terminate.
func (b *Bot) Stop() {
	b.lifeMu.Lock()
	if b.lifeStop != nil {
		b.lifeStop()
	}
	b.lifeMu.Unlock()
}

// lifecycleContext returns the context of the current run, or
// context.Background() when the bot is idle.
func (b *Bot) lifecycleContext() context.Context {
	b.lifeMu.Lock()
	defer b.lifeMu.Unlock()
	if b.lifeCtx != nil {
		return b.lifeCtx
	}
	return context.Background()
}

func parseIntDefault(s string, def int) int {
	if n, err := strconv.Atoi(s); err == nil {
		return n
	}
	return def
}

func (b *Bot) onStart(ctx context.Context) error {
	b.mu.RLock()
	hooks := append([]hookEntry(nil), b.startHandlers...)
	b.mu.RUnlock()
	for _, h := range hooks {
		if err := h.fn(ctx, b); err != nil {
			return err
		}
	}
	return nil
}

func (b *Bot) onStop(ctx context.Context) {
	b.mu.RLock()
	hooks := append([]hookEntry(nil), b.stopHandlers...)
	b.mu.RUnlock()
	for i := len(hooks) - 1; i >= 0; i-- {
		func() {
			defer func() {
				if r := recover(); r != nil {
					b.errorf("stop hook panic: %v", r)
				}
			}()
			if err := hooks[i].fn(ctx, b); err != nil {
				b.errorf("stop hook: %v", err)
			}
		}()
	}
}

// drain waits for in-flight handlers, bounded by ShutdownTimeout.
func (b *Bot) drain() {
	if !b.disp.wait(b.opts.ShutdownTimeout) {
		b.warnf("shutdown timeout, %d chats still busy", b.disp.pending())
	}
}

// shortHash returns a 16-char hex digest of s (SHA-256 truncated).
func shortHash(s string) string {
	h := sha256Sum([]byte(s))
	return fmt.Sprintf("%x", h[:8])
}

// jsonEncode is a small helper to avoid importing encoding/json in every file.
func jsonEncode(v any) ([]byte, error) { return json.Marshal(v) }

// Constants re-exported so callers don't have to import core.
const (
	UserAgent = core.UserAgent
	Version2  = core.Version
)
