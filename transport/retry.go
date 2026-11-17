package transport

import (
	"context"
	"errors"
	"math"
	"math/rand/v2"
	"time"
)

// RetryPolicy describes how a failed request should be retried.
type RetryPolicy struct {
	MaxRetries int
	BaseDelay  time.Duration
	MaxDelay   time.Duration
	RetryOn    func(status int, err error) bool
}

// DefaultRetryPolicy returns the retry policy used by ParsRubika by
// default for idempotent requests (GET): 5 attempts, 1-second base
// delay, capped at 30 seconds, and retry on transient network errors
// and 5xx responses. Because GET is idempotent, retrying on an
// ambiguous transport error (the server may or may not have seen the
// request) is safe.
func DefaultRetryPolicy() RetryPolicy {
	return RetryPolicy{
		MaxRetries: 5,
		BaseDelay:  time.Second,
		MaxDelay:   30 * time.Second,
		RetryOn: func(status int, err error) bool {
			if err != nil {
				return true
			}
			return status == 0 || status >= 500 || status == 408 || status == 425 || status == 429
		},
	}
}

// PostRetryPolicy is the retry policy for non-idempotent POST requests
// (every Rubika message-send, file-upload-start, etc.). It retries ONLY
// on explicit, server-signalled retryable HTTP statuses (408, 425, 429,
// 500, 502, 503, 504) and NEVER on a transport-level error
// (status 0 / err != nil).
//
// A transport error on a POST is ambiguous: the request may have
// reached the server and been applied (e.g. the message may already be
// delivered) even though the client never saw the response. Retrying
// in that case risks delivering a DUPLICATE message, so the request is
// failed back to the caller instead. The server-signalled statuses are
// safe to retry because the server explicitly states it did not apply
// the request (or is asking the client to slow down).
func PostRetryPolicy() RetryPolicy {
	return RetryPolicy{
		MaxRetries: 5,
		BaseDelay:  time.Second,
		MaxDelay:   30 * time.Second,
		RetryOn: func(status int, err error) bool {
			if err != nil {
				return false
			}
			switch status {
			case 408, 425, 429, 500, 502, 503, 504:
				return true
			}
			return false
		},
	}
}

// Sleep waits for the exponential backoff of the supplied attempt
// number, with up to ±25 % jitter so that a burst of retries from many
// clients does not synchronise (thundering herd). The first attempt
// uses base delay; subsequent attempts double the delay each time,
// capped by `MaxDelay`.
func (p RetryPolicy) Sleep(ctx context.Context, attempt int) error {
	if p.BaseDelay <= 0 || attempt <= 0 {
		return nil
	}
	delay := time.Duration(math.Pow(2, float64(attempt-1))) * p.BaseDelay
	if p.MaxDelay > 0 && delay > p.MaxDelay {
		delay = p.MaxDelay
	}
	// ±25 % full jitter to desynchronise concurrent retry storms.
	if delay > 0 {
		jitter := time.Duration(rand.Int64N(int64(delay) / 2))
		if attempt&1 == 0 {
			delay += jitter
		} else {
			delay -= jitter
		}
	}
	select {
	case <-time.After(delay):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// ShouldRetry reports whether the policy says the request should be
// retried for the given (status, err) pair.
func (p RetryPolicy) ShouldRetry(status int, err error) bool {
	if p.RetryOn == nil {
		return err != nil
	}
	return p.RetryOn(status, err)
}

// ErrMaxRetries is returned when a request exhausts all retries.
var ErrMaxRetries = errors.New("transport: max retries exceeded")
