package bot

import (
	"context"
	"sync"
	"time"
)

// ===== EventBus =======================================================
//
// EventBus is a simple in-memory pub/sub system. Bots can use it to
// decouple event producers from consumers (e.g. a message handler that
// emits a "new_user" event, and a separate welcome module that listens).

// EventHandler is a callback invoked when a matching event is published.
type EventHandler func(ctx context.Context, data any)

// EventBus is a topic-based in-memory event bus. It is safe for
// concurrent use. Subscriptions are stored in a map keyed by topic.
type EventBus struct {
	mu          sync.RWMutex
	subscribers map[string]map[int]EventHandler
	nextID      int
}

// NewEventBus creates a new EventBus.
func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[string]map[int]EventHandler),
	}
}

// Subscribe registers a handler for the given topic. Returns an
// unsubscribe function.
func (eb *EventBus) Subscribe(topic string, handler EventHandler) func() {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	if _, ok := eb.subscribers[topic]; !ok {
		eb.subscribers[topic] = make(map[int]EventHandler)
	}
	id := eb.nextID
	eb.nextID++
	eb.subscribers[topic][id] = handler
	return func() {
		eb.mu.Lock()
		defer eb.mu.Unlock()
		if subs, ok := eb.subscribers[topic]; ok {
			delete(subs, id)
			if len(subs) == 0 {
				delete(eb.subscribers, topic)
			}
		}
	}
}

// Publish emits an event to all subscribers of the topic. Each handler
// is invoked in its own goroutine so a slow handler does not block
// others. The context is passed to each handler.
func (eb *EventBus) Publish(ctx context.Context, topic string, data any) {
	eb.mu.RLock()
	subs := make([]EventHandler, 0, len(eb.subscribers[topic]))
	for _, h := range eb.subscribers[topic] {
		subs = append(subs, h)
	}
	eb.mu.RUnlock()
	for _, h := range subs {
		go h(ctx, data)
	}
}

// PublishSync emits an event synchronously — handlers are called in
// sequence in the calling goroutine. A panicking handler is recovered
// so the bus is never poisoned.
func (eb *EventBus) PublishSync(ctx context.Context, topic string, data any) {
	eb.mu.RLock()
	subs := make([]EventHandler, 0, len(eb.subscribers[topic]))
	for _, h := range eb.subscribers[topic] {
		subs = append(subs, h)
	}
	eb.mu.RUnlock()
	for _, h := range subs {
		func() {
			defer func() { _ = recover() }()
			h(ctx, data)
		}()
	}
}

// SubscriberCount returns the number of subscribers for a topic.
func (eb *EventBus) SubscriberCount(topic string) int {
	eb.mu.RLock()
	defer eb.mu.RUnlock()
	return len(eb.subscribers[topic])
}

// Topics returns the list of active topics.
func (eb *EventBus) Topics() []string {
	eb.mu.RLock()
	defer eb.mu.RUnlock()
	topics := make([]string, 0, len(eb.subscribers))
	for t := range eb.subscribers {
		topics = append(topics, t)
	}
	return topics
}

// ===== TokenBucket rate limiter =======================================
//
// TokenBucket is a classic token-bucket rate limiter. It is safe for
// concurrent use. Tokens are refilled at a fixed rate up to capacity.

// TokenBucket implements a token-bucket rate limiter.
type TokenBucket struct {
	mu         sync.Mutex
	capacity   float64
	tokens     float64
	fillRate   float64 // tokens per second
	lastRefill time.Time
}

// NewTokenBucket creates a TokenBucket with the given capacity (max
// tokens) and fillRate (tokens added per second).
func NewTokenBucket(capacity float64, fillRate float64) *TokenBucket {
	return &TokenBucket{
		capacity:   capacity,
		tokens:     capacity,
		fillRate:   fillRate,
		lastRefill: time.Now(),
	}
}

// Allow attempts to consume 1 token. Returns true if a token was
// available (and consumed), false otherwise.
func (tb *TokenBucket) Allow() bool {
	return tb.AllowN(1)
}

// AllowN attempts to consume n tokens. Returns true if n tokens were
// available (and consumed), false otherwise.
func (tb *TokenBucket) AllowN(n float64) bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	tb.refill()
	if tb.tokens >= n {
		tb.tokens -= n
		return true
	}
	return false
}

// Wait blocks until a token is available (or ctx is cancelled). Returns
// nil on success, ctx.Err() on cancellation.
func (tb *TokenBucket) Wait(ctx context.Context) error {
	for {
		tb.mu.Lock()
		tb.refill()
		if tb.tokens >= 1 {
			tb.tokens -= 1
			tb.mu.Unlock()
			return nil
		}
		// Calculate wait time
		deficit := 1 - tb.tokens
		wait := time.Duration(deficit/tb.fillRate*float64(time.Second)) + 1
		tb.mu.Unlock()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(wait):
		}
	}
}

// refill adds tokens based on elapsed time. Must be called under lock.
func (tb *TokenBucket) refill() {
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()
	tb.tokens += elapsed * tb.fillRate
	if tb.tokens > tb.capacity {
		tb.tokens = tb.capacity
	}
	tb.lastRefill = now
}

// Tokens returns the current number of available tokens (after refill).
func (tb *TokenBucket) Tokens() float64 {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	tb.refill()
	return tb.tokens
}

// ===== MetricsTracker =================================================
//
// MetricsTracker is a lightweight in-memory metrics collector. It
// supports counters (monotonic) and gauges (settable). Useful for
// bot analytics without an external metrics backend.

// MetricsTracker holds named counters and gauges in memory.
type MetricsTracker struct {
	mu       sync.RWMutex
	counters map[string]float64
	gauges   map[string]float64
}

// NewMetricsTracker creates a MetricsTracker.
func NewMetricsTracker() *MetricsTracker {
	return &MetricsTracker{
		counters: make(map[string]float64),
		gauges:   make(map[string]float64),
	}
}

// IncrementCounter adds value to the named counter.
func (m *MetricsTracker) IncrementCounter(name string, value float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.counters[name] += value
}

// SetGauge sets the named gauge to value.
func (m *MetricsTracker) SetGauge(name string, value float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.gauges[name] = value
}

// GetCounter returns the current counter value.
func (m *MetricsTracker) GetCounter(name string) float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.counters[name]
}

// GetGauge returns the current gauge value.
func (m *MetricsTracker) GetGauge(name string) float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.gauges[name]
}

// AllMetrics returns snapshots of all counters and gauges.
func (m *MetricsTracker) AllMetrics() (counters, gauges map[string]float64) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	c := make(map[string]float64, len(m.counters))
	for k, v := range m.counters {
		c[k] = v
	}
	g := make(map[string]float64, len(m.gauges))
	for k, v := range m.gauges {
		g[k] = v
	}
	return c, g
}

// Reset clears all counters and gauges.
func (m *MetricsTracker) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.counters = make(map[string]float64)
	m.gauges = make(map[string]float64)
}
