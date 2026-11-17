// Package guard provides safety guards — rate limiting, timeout
// enforcement, and circuit-breaker patterns that protect bots from
// runaway loops, accidental spam, and server-side bans.
package guard

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Circuit is a simple circuit-breaker: when the failure count
// exceeds `threshold` in `window`, the circuit opens for `cooldown`.
// After `cooldown`, the circuit enters the half-open state and
// allows ONE trial; further Allow() calls in half-open state are
// rejected until that trial resolves via Success() or Failure().
// If the trial succeeds, the circuit closes; if it fails, the
// circuit re-opens for another cooldown period.
type Circuit struct {
	mu                    sync.Mutex
	threshold             int
	window                time.Duration
	cooldown              time.Duration
	failures              []time.Time
	opened                time.Time
	state                 circuitState
	halfOpenTrialInFlight bool
}

type circuitState int

const (
	circuitClosed circuitState = iota
	circuitOpen
	circuitHalfOpen
)

// NewCircuit constructs a circuit breaker.
func NewCircuit(threshold int, window, cooldown time.Duration) *Circuit {
	return &Circuit{threshold: threshold, window: window, cooldown: cooldown, state: circuitClosed}
}

// Allow reports whether a new request should be permitted.
func (c *Circuit) Allow() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	switch c.state {
	case circuitOpen:
		if now.Sub(c.opened) >= c.cooldown {
			c.state = circuitHalfOpen
			c.halfOpenTrialInFlight = true
			return true
		}
		return false
	case circuitHalfOpen:
		// Only one trial is allowed in half-open state. Subsequent
		// Allow() calls are rejected until the trial resolves.
		if c.halfOpenTrialInFlight {
			return false
		}
		c.halfOpenTrialInFlight = true
		return true
	default: // circuitClosed
		return true
	}
}

// Success records a successful call.
func (c *Circuit) Success() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.state = circuitClosed
	c.failures = nil
	c.halfOpenTrialInFlight = false
}

// Failure records a failed call.
func (c *Circuit) Failure() {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	// In half-open state, a single failure re-opens the circuit.
	if c.state == circuitHalfOpen {
		c.state = circuitOpen
		c.opened = now
		c.halfOpenTrialInFlight = false
		return
	}
	c.halfOpenTrialInFlight = false
	cutoff := now.Add(-c.window)
	// Drop expired failures.
	var fresh []time.Time
	for _, t := range c.failures {
		if t.After(cutoff) {
			fresh = append(fresh, t)
		}
	}
	fresh = append(fresh, now)
	c.failures = fresh
	if len(fresh) >= c.threshold {
		c.state = circuitOpen
		c.opened = now
	}
}

// Timeout wraps a context with the supplied timeout. It is a
// convenience wrapper for context.WithTimeout.
func Timeout(parent context.Context, d time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, d)
}

// Retry retries `fn` up to `max` times with exponential backoff.
// The first retry waits `base`; each subsequent retry doubles the
// delay. Returns the last error when all retries are exhausted.
func Retry(ctx context.Context, max int, base time.Duration, fn func() error) error {
	// Guard against a misconfigured caller: with max < 1 the loop body
	// never runs and lastErr stays nil, which would previously report
	// success without ever calling fn. Surface it as a real error
	// instead of a silent nil.
	if max < 1 {
		return fmt.Errorf("guard: Retry requires max >= 1 (got %d)", max)
	}
	if fn == nil {
		return fmt.Errorf("guard: Retry requires a non-nil fn")
	}
	var lastErr error
	delay := base
	for attempt := 0; attempt < max; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		lastErr = fn()
		if lastErr == nil {
			return nil
		}
		if attempt < max-1 {
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return ctx.Err()
			}
			delay *= 2
		}
	}
	return lastErr
}
