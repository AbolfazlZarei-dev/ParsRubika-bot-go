# `transport` — Ratelimit

_Source: [`ratelimit.go`](../../ParsRubika/transport/ratelimit.go)_

## Symbols

### `RateLimiter`

```go
type RateLimiter struct
```

RateLimiter enforces a minimum delay between successive requests.
It is a simple token-bucket: the bucket holds exactly one token,
which is consumed by each request and refilled after `interval`.

The implementation is correct under concurrency: each caller reserves
a projected issue time (lastIssued + interval) while holding the
lock, so concurrent callers are serialised rather than all sleeping
the same duration and proceeding simultaneously.

### `NewRateLimiter`

```go
func NewRateLimiter(interval time.Duration) *RateLimiter
```

NewRateLimiter returns a RateLimiter that waits at least `interval`
between requests.

### `Wait`

```go
func (r *RateLimiter) Wait(ctx context.Context) error
```

Wait blocks until enough time has elapsed since the previous Wait
call. Cancellation via `ctx` is honoured.

Under concurrency, each caller reserves a future issue slot
(lastIssued + interval) atomically, so callers are effectively
queued rather than racing. The projected issue time is recorded
immediately so the next caller computes its wait relative to this
caller's slot, not the wall clock.
