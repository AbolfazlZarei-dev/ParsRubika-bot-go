# `bot` — Eventbus

_Source: [`eventbus.go`](../../ParsRubika/bot/eventbus.go)_

## Symbols

### `EventHandler`

```go
type EventHandler func(ctx context.Context, data any)
```

EventHandler is a callback invoked when a matching event is published.

### `EventBus`

```go
type EventBus struct
```

EventBus is a topic-based in-memory event bus. It is safe for
concurrent use. Subscriptions are stored in a map keyed by topic.

### `NewEventBus`

```go
func NewEventBus() *EventBus
```

NewEventBus creates a new EventBus.

### `Subscribe`

```go
func (eb *EventBus) Subscribe(topic string, handler EventHandler) func()
```

Subscribe registers a handler for the given topic. Returns an
unsubscribe function.

### `Publish`

```go
func (eb *EventBus) Publish(ctx context.Context, topic string, data any)
```

Publish emits an event to all subscribers of the topic. Each handler
is invoked in its own goroutine so a slow handler does not block
others. The context is passed to each handler.

### `PublishSync`

```go
func (eb *EventBus) PublishSync(ctx context.Context, topic string, data any)
```

PublishSync emits an event synchronously — handlers are called in
sequence in the calling goroutine. A panicking handler is recovered
so the bus is never poisoned.

### `SubscriberCount`

```go
func (eb *EventBus) SubscriberCount(topic string) int
```

SubscriberCount returns the number of subscribers for a topic.

### `Topics`

```go
func (eb *EventBus) Topics() []string
```

Topics returns the list of active topics.

### `TokenBucket`

```go
type TokenBucket struct
```

TokenBucket implements a token-bucket rate limiter.

### `NewTokenBucket`

```go
func NewTokenBucket(capacity float64, fillRate float64) *TokenBucket
```

NewTokenBucket creates a TokenBucket with the given capacity (max
tokens) and fillRate (tokens added per second).

### `Allow`

```go
func (tb *TokenBucket) Allow() bool
```

Allow attempts to consume 1 token. Returns true if a token was
available (and consumed), false otherwise.

### `AllowN`

```go
func (tb *TokenBucket) AllowN(n float64) bool
```

AllowN attempts to consume n tokens. Returns true if n tokens were
available (and consumed), false otherwise.

### `Wait`

```go
func (tb *TokenBucket) Wait(ctx context.Context) error
```

Wait blocks until a token is available (or ctx is cancelled). Returns
nil on success, ctx.Err() on cancellation.

### `Tokens`

```go
func (tb *TokenBucket) Tokens() float64
```

Tokens returns the current number of available tokens (after refill).

### `MetricsTracker`

```go
type MetricsTracker struct
```

MetricsTracker holds named counters and gauges in memory.

### `NewMetricsTracker`

```go
func NewMetricsTracker() *MetricsTracker
```

NewMetricsTracker creates a MetricsTracker.

### `IncrementCounter`

```go
func (m *MetricsTracker) IncrementCounter(name string, value float64)
```

IncrementCounter adds value to the named counter.

### `SetGauge`

```go
func (m *MetricsTracker) SetGauge(name string, value float64)
```

SetGauge sets the named gauge to value.

### `GetCounter`

```go
func (m *MetricsTracker) GetCounter(name string) float64
```

GetCounter returns the current counter value.

### `GetGauge`

```go
func (m *MetricsTracker) GetGauge(name string) float64
```

GetGauge returns the current gauge value.

### `AllMetrics`

```go
func (m *MetricsTracker) AllMetrics() (counters, gauges map[string]float64)
```

AllMetrics returns snapshots of all counters and gauges.

### `Reset`

```go
func (m *MetricsTracker) Reset()
```

Reset clears all counters and gauges.
