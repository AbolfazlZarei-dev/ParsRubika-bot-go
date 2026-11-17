package bot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/transport"
)

// retryableAPIStatuses is the set of Rubika API `status` strings that
// indicate a transient, retryable failure.
var retryableAPIStatuses = map[string]struct{}{
	"TooRequests":        {},
	"ServerError":        {},
	"503":                {},
	"504":                {},
	"BadGateway":         {},
	"ServiceUnavailable": {},
}

// maxResponseBytes bounds the size of an API response held in memory.
const maxResponseBytes = 32 << 20

// makeRequest performs a rate-limited API call with retries.
func (b *Bot) makeRequest(ctx context.Context, method string, payload any) (map[string]any, error) {
	return b.call(ctx, method, payload, true)
}

// call is the shared request engine. When limited is false the global
// rate limiter is skipped (used by the polling loop, which must never
// compete with message sending for slots).
func (b *Bot) call(ctx context.Context, method string, payload any, limited bool) (map[string]any, error) {
	if limited {
		if err := b.rateLimiter.Wait(ctx); err != nil {
			return nil, err
		}
	}
	url := b.baseURL + method
	buf := transport.AcquireBuffer()
	defer transport.ReleaseBuffer(buf)
	if err := json.NewEncoder(buf).Encode(payload); err != nil {
		return nil, fmt.Errorf("bot: marshal payload: %w", err)
	}
	body := buf.Bytes()
	// json.Encoder.Encode appends a trailing newline; trim it so the
	// body matches json.Marshal output and Content-Length stays exact.
	if n := len(body); n > 0 && body[n-1] == '\n' {
		body = body[:n-1]
	}
	b.stats.APICalls.Add(1)

	retries := b.opts.MaxRetries
	if retries < 0 {
		retries = 0
	}
	var lastErr error
	for attempt := 0; attempt <= retries; attempt++ {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		data, retryAfter, err := b.attempt(ctx, url, method, body)
		if err == nil {
			return data, nil
		}
		lastErr = err
		if attempt == retries || !b.shouldRetry(err) {
			break
		}
		b.warnf("%s attempt %d/%d failed: %v", method, attempt+1, retries+1, err)
		if !b.sleepBackoff(ctx, attempt, retryAfter) {
			return nil, ctx.Err()
		}
	}
	if ae, ok := lastErr.(*APIError); ok {
		return nil, ae
	}
	return nil, fmt.Errorf("bot: %s failed: %w", method, lastErr)
}

// attempt performs one HTTP round trip.
func (b *Bot) attempt(ctx context.Context, url, method string, body []byte) (map[string]any, time.Duration, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", b.opts.UserAgent)
	resp, err := b.transport.HTTPClient().Do(req)
	if err != nil {
		return nil, 0, err
	}
	raw, readErr := io.ReadAll(io.LimitReader(resp.Body, maxResponseBytes))
	_ = resp.Body.Close()
	retryAfter := parseRetryAfter(resp.Header.Get("Retry-After"))
	if readErr != nil {
		return nil, retryAfter, fmt.Errorf("read response: %w", readErr)
	}
	if resp.StatusCode >= 400 {
		return nil, retryAfter, &APIError{
			Status:     strconv.Itoa(resp.StatusCode),
			DevMessage: string(raw),
			HTTPStatus: resp.StatusCode,
			RetryAfter: retryAfter,
		}
	}
	var parsed map[string]any
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, 0, &APIError{Status: "InvalidResponse", DevMessage: string(raw)}
	}
	status, _ := parsed["status"].(string)
	if status != "OK" {
		return nil, retryAfter, &APIError{
			Status:     status,
			DevMessage: parsedString(parsed, "dev_message"),
			HTTPStatus: resp.StatusCode,
			RetryAfter: retryAfter,
		}
	}
	data, ok := parsed["data"].(map[string]any)
	if !ok {
		// Some methods legitimately answer without a data object.
		return map[string]any{}, 0, nil
	}
	return data, 0, nil
}

// shouldRetry classifies an error as transient.
func (b *Bot) shouldRetry(err error) bool {
	if ae, ok := err.(*APIError); ok {
		if ae.Status == "InvalidResponse" {
			return false
		}
		if ae.HTTPStatus != 0 {
			if _, ok := b.opts.RetryStatuses[ae.HTTPStatus]; ok {
				return true
			}
			return false
		}
		return isRetryableAPIStatus(ae.Status)
	}
	// Network level failures (timeouts, resets, truncated bodies).
	return true
}

func parsedString(m map[string]any, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// isRetryableAPIStatus reports whether a Rubika API `status` string
// names a transient error worth retrying.
func isRetryableAPIStatus(status string) bool {
	_, ok := retryableAPIStatuses[status]
	return ok
}

// parseRetryAfter understands both the delta-seconds and HTTP-date forms.
func parseRetryAfter(v string) time.Duration {
	if v == "" {
		return 0
	}
	if secs, err := strconv.Atoi(v); err == nil && secs >= 0 {
		return time.Duration(secs) * time.Second
	}
	if t, err := http.ParseTime(v); err == nil {
		if d := time.Until(t); d > 0 {
			return d
		}
	}
	return 0
}

// backoffDelay computes the exponential delay for an attempt with
// "equal jitter", capped at MaxBackoff.
func (b *Bot) backoffDelay(attempt int) time.Duration {
	base := b.opts.BackoffFactor
	if base <= 0 {
		return 0
	}
	if attempt > 20 {
		attempt = 20
	}
	d := base << uint(attempt)
	if max := b.opts.MaxBackoff; max > 0 && (d > max || d <= 0) {
		d = max
	}
	half := d / 2
	if half <= 0 {
		return d
	}
	return half + time.Duration(rand.Int63n(int64(half)+1))
}

// sleepBackoff waits for the larger of the computed backoff and the
// server supplied Retry-After. It returns false when ctx ended first.
func (b *Bot) sleepBackoff(ctx context.Context, attempt int, retryAfter time.Duration) bool {
	d := b.backoffDelay(attempt)
	if retryAfter > d {
		d = retryAfter
	}
	if d <= 0 {
		return ctx.Err() == nil
	}
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-t.C:
		return true
	case <-ctx.Done():
		return false
	}
}
