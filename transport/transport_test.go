package transport

import (
	"context"
	"testing"
	"time"
)

func TestRateLimiter(t *testing.T) {
	rl := NewRateLimiter(10 * time.Millisecond)
	ctx := context.Background()
	// First call returns immediately.
	start := time.Now()
	_ = rl.Wait(ctx)
	if elapsed := time.Since(start); elapsed > 5*time.Millisecond {
		t.Errorf("first Wait took too long: %v", elapsed)
	}
	// Second call waits at least 10ms.
	start = time.Now()
	_ = rl.Wait(ctx)
	if elapsed := time.Since(start); elapsed < 8*time.Millisecond {
		t.Errorf("second Wait did not rate-limit: %v", elapsed)
	}
}

func TestRateLimiterZero(t *testing.T) {
	rl := NewRateLimiter(0)
	ctx := context.Background()
	start := time.Now()
	_ = rl.Wait(ctx)
	_ = rl.Wait(ctx)
	if elapsed := time.Since(start); elapsed > 5*time.Millisecond {
		t.Errorf("zero-interval rate limiter should not wait: %v", elapsed)
	}
}

func TestRetryPolicySleep(t *testing.T) {
	p := RetryPolicy{BaseDelay: 1 * time.Millisecond, MaxDelay: 10 * time.Millisecond}
	ctx := context.Background()
	start := time.Now()
	_ = p.Sleep(ctx, 1)
	if elapsed := time.Since(start); elapsed < 1*time.Millisecond {
		t.Errorf("Sleep(1) took only %v", elapsed)
	}
	// Cancellation propagates.
	cancelCtx, cancel := context.WithCancel(ctx)
	cancel()
	if err := p.Sleep(cancelCtx, 1); err == nil {
		t.Error("Sleep should return context error after cancellation")
	}
}

func TestRetryPolicyShouldRetry(t *testing.T) {
	p := DefaultRetryPolicy()
	if !p.ShouldRetry(0, errDummy{}) {
		t.Error("ShouldRetry should retry on errors")
	}
	if !p.ShouldRetry(500, nil) {
		t.Error("ShouldRetry should retry on 5xx")
	}
	if p.ShouldRetry(200, nil) {
		t.Error("ShouldRetry should NOT retry on 200")
	}
	if !p.ShouldRetry(429, nil) {
		t.Error("ShouldRetry should retry on 429")
	}
}

type errDummy struct{}

func (errDummy) Error() string { return "dummy" }
