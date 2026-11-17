# `guard` — Package guard provides safety guards — rate limiting, timeout

_Source: [`guard.go`](../../../ParsRubika/core/guard/guard.go)_

## Symbols

### `Circuit`

```go
type Circuit struct
```

Circuit is a simple circuit-breaker: when the failure count
exceeds `threshold` in `window`, the circuit opens for `cooldown`.
After `cooldown`, the circuit enters the half-open state and
allows ONE trial; further Allow() calls in half-open state are
rejected until that trial resolves via Success() or Failure().
If the trial succeeds, the circuit closes; if it fails, the
circuit re-opens for another cooldown period.

### `NewCircuit`

```go
func NewCircuit(threshold int, window, cooldown time.Duration) *Circuit
```

NewCircuit constructs a circuit breaker.

### `Allow`

```go
func (c *Circuit) Allow() bool
```

Allow reports whether a new request should be permitted.

### `Success`

```go
func (c *Circuit) Success()
```

Success records a successful call.

### `Failure`

```go
func (c *Circuit) Failure()
```

Failure records a failed call.

### `Timeout`

```go
func Timeout(parent context.Context, d time.Duration) (context.Context, context.CancelFunc)
```

Timeout wraps a context with the supplied timeout. It is a
convenience wrapper for context.WithTimeout.

### `Retry`

```go
func Retry(ctx context.Context, max int, base time.Duration, fn func() error) error
```

Retry retries `fn` up to `max` times with exponential backoff.
The first retry waits `base`; each subsequent retry doubles the
delay. Returns the last error when all retries are exhausted.
