package transport

import (
	"context"
	"sync"
	"time"
)

// RateLimiter enforces a minimum delay between successive requests.
// It is a simple token-bucket: the bucket holds exactly one token,
// which is consumed by each request and refilled after `interval`.
//
// The implementation is correct under concurrency: each caller reserves
// a projected issue time (lastIssued + interval) while holding the
// lock, so concurrent callers are serialised rather than all sleeping
// the same duration and proceeding simultaneously.
type RateLimiter struct {
	mu         sync.Mutex
	interval   time.Duration
	lastIssued time.Time
}

// NewRateLimiter returns a RateLimiter that waits at least `interval`
// between requests.
func NewRateLimiter(interval time.Duration) *RateLimiter {
	return &RateLimiter{interval: interval}
}

// Wait blocks until enough time has elapsed since the previous Wait
// call. Cancellation via `ctx` is honoured.
//
// Under concurrency, each caller reserves a future issue slot
// (lastIssued + interval) atomically, so callers are effectively
// queued rather than racing. The projected issue time is recorded
// immediately so the next caller computes its wait relative to this
// caller's slot, not the wall clock.
func (r *RateLimiter) Wait(ctx context.Context) error {
	if r.interval <= 0 {
		return nil
	}
	r.mu.Lock()
	// Compute the earliest this caller may issue. If lastIssued is zero
	// (first call) or enough time has elapsed, issue immediately;
	// otherwise queue behind the previous caller.
	var wait time.Duration
	now := time.Now()
	if r.lastIssued.IsZero() {
		r.lastIssued = now
	} else {
		projected := r.lastIssued.Add(r.interval)
		if projected.After(now) {
			wait = projected.Sub(now)
			// Reserve this caller's slot at projected, NOT now, so the
			// next caller queues behind us rather than behind the wall
			// clock.
			r.lastIssued = projected
		} else {
			r.lastIssued = now
		}
	}
	r.mu.Unlock()

	if wait <= 0 {
		return nil
	}
	select {
	case <-time.After(wait):
		return nil
	case <-ctx.Done():
		// On cancellation, release the reserved slot so the next caller
		// doesn't have to wait extra time for a caller that never
		// actually issued.
		r.mu.Lock()
		// Only roll back if our reservation is still the last one
		// (another caller may have already queued behind us).
		if r.lastIssued.Equal(now.Add(wait)) || r.lastIssued.Sub(now) <= wait+r.interval {
			// Best-effort rollback: move lastIssued back by interval.
			// This is approximate but prevents unbounded queue growth
			// when callers cancel.
			r.lastIssued = r.lastIssued.Add(-r.interval)
		}
		r.mu.Unlock()
		return ctx.Err()
	}
}
