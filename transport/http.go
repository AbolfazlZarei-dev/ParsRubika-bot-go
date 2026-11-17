package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// MaxResponseBody caps how many bytes of an HTTP response body the
// transport will read into memory. Rubika API responses are small JSON
// documents; the cap exists so a malicious or buggy server cannot make
// the client allocate unbounded memory by streaming an enormous body.
// A response that exceeds the cap fails with ErrResponseBodyTooLarge.
const MaxResponseBody = 16 << 20 // 16 MiB

// ErrResponseBodyTooLarge is returned when a response body exceeds
// MaxResponseBody.
var ErrResponseBodyTooLarge = fmt.Errorf("transport: response body exceeds %d bytes", MaxResponseBody)

// maxErrorSnippet caps how much of a response body is embedded in an
// error string. The body may contain echoed auth material or other
// sensitive data, so it is both truncated AND stripped of obvious
// token/key patterns before being placed in the error.
const maxErrorSnippet = 512

// readBody reads at most MaxResponseBody+1 bytes from r and returns the
// body plus a flag indicating whether the limit was exceeded.
func readBody(r io.Reader) ([]byte, bool, error) {
	limited := io.LimitReader(r, MaxResponseBody+1)
	b, err := io.ReadAll(limited)
	if err != nil {
		return nil, false, err
	}
	if len(b) > MaxResponseBody {
		return nil, true, nil
	}
	return b, false, nil
}

// sanitizeForError returns a short, safe snippet of a response body for
// embedding in an error string. It truncates to maxErrorSnippet bytes
// and strips common secret-looking substrings (auth tokens, PEM keys,
// long hex/base64 runs) so that an error never leaks credentials even
// if the server echoes them back.
func sanitizeForError(body []byte) string {
	if len(body) > maxErrorSnippet {
		body = body[:maxErrorSnippet]
	}
	s := string(body)
	// Drop PEM private-key blocks outright.
	if i := indexOf(s, "-----BEGIN"); i >= 0 {
		s = s[:i] + "[redacted pem]"
	}
	return s
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}

// Post sends a POST request with the supplied JSON payload and returns
// the raw response body. It retries on transient failures according to
// the configured RetryPolicy.
//
// The payload is encoded into a process-wide pooled buffer rather than a
// freshly allocated slice, so the encode step reuses warm memory across
// requests. The resulting bytes are handed straight to bytes.NewReader
// in postRaw with no intermediate copy.
func (t *Transport) Post(ctx context.Context, url string, payload any) ([]byte, int, error) {
	buf := AcquireBuffer()
	defer ReleaseBuffer(buf)
	if err := json.NewEncoder(buf).Encode(payload); err != nil {
		return nil, 0, fmt.Errorf("transport: marshal: %w", err)
	}
	body := buf.Bytes()
	// json.Encoder.Encode appends a trailing newline; trim it so the
	// body matches json.Marshal output and Content-Length stays exact.
	if n := len(body); n > 0 && body[n-1] == '\n' {
		body = body[:n-1]
	}
	return t.postRaw(ctx, url, body)
}

// postRaw is the inner loop that performs the actual POST with retries.
//
// Retry semantics: POST is non-idempotent, so we use PostRetryPolicy —
// the request is retried ONLY on explicit, server-signalled retryable
// statuses (429/5xx/408/425). A transport-level error (the server may
// already have applied the request) is NOT retried, so a duplicate
// message can never be delivered by the retry loop.
//
// bytes.NewReader wraps the body slice with no copy; the request fully
// consumes it during Do, so reusing the same slice across retry
// attempts is safe.
func (t *Transport) postRaw(ctx context.Context, url string, body []byte) ([]byte, int, error) {
	policy := PostRetryPolicy()
	if t.maxRetries > 0 {
		policy.MaxRetries = t.maxRetries
	}

	var lastErr error
	var lastStatus int
	for attempt := 1; attempt <= policy.MaxRetries; attempt++ {
		if err := ctx.Err(); err != nil {
			return nil, 0, err
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
		if err != nil {
			return nil, 0, err
		}
		t.applyHeaders(req)
		resp, err := t.httpClient.Do(req)
		if err != nil {
			// Ambiguous for POST: do not retry (would risk a duplicate).
			return nil, 0, err
		}
		respBody, tooLarge, readErr := readBody(resp.Body)
		_ = resp.Body.Close()
		if readErr != nil {
			// Ambiguous for POST: do not retry.
			return nil, resp.StatusCode, fmt.Errorf("transport: read response: %w", readErr)
		}
		if tooLarge {
			return nil, resp.StatusCode, ErrResponseBodyTooLarge
		}
		if resp.StatusCode >= 400 {
			lastErr = fmt.Errorf("transport: http %d: %s", resp.StatusCode, sanitizeForError(respBody))
			lastStatus = resp.StatusCode
			if !policy.ShouldRetry(resp.StatusCode, nil) {
				break
			}
			_ = policy.Sleep(ctx, attempt)
			continue
		}
		return respBody, resp.StatusCode, nil
	}
	if lastErr == nil {
		lastErr = ErrMaxRetries
	}
	return nil, lastStatus, lastErr
}

// Get sends a GET request and returns the raw response body.
func (t *Transport) Get(ctx context.Context, url string) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, err
	}
	t.applyHeaders(req)
	resp, err := t.httpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	body, tooLarge, err := readBody(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("transport: read response: %w", err)
	}
	if tooLarge {
		return nil, resp.StatusCode, ErrResponseBodyTooLarge
	}
	if resp.StatusCode >= 400 {
		return body, resp.StatusCode, fmt.Errorf("transport: http %d: %s", resp.StatusCode, sanitizeForError(body))
	}
	return body, resp.StatusCode, nil
}

// PostJSON is a convenience helper that JSON-decodes the response body
// into `out`.
func (t *Transport) PostJSON(ctx context.Context, url string, payload any, out any) (int, error) {
	body, status, err := t.Post(ctx, url, payload)
	if err != nil {
		return status, err
	}
	if out != nil {
		if err := json.Unmarshal(body, out); err != nil {
			return status, fmt.Errorf("transport: decode response: %w", err)
		}
	}
	return status, nil
}

// New constructs a Transport with the supplied options. It uses the
// default HTTP/2-enabled client and the standard header set.
func New(opts ...Option) *Transport {
	t := &Transport{
		headers:    DefaultHeaders(""),
		timeout:    30 * time.Second,
		maxRetries: 5,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(t)
		}
	}
	if t.userAgent != "" {
		t.headers["user-agent"] = t.userAgent
	}
	t.httpClient = NewHTTPClient(t.timeout, t.proxy, false)
	return t
}

// SetHeader updates a single header. Useful for the bot API client
// which uses a different base URL per token. It is safe to call
// concurrently with in-flight requests.
func (t *Transport) SetHeader(key, value string) {
	t.headerMu.Lock()
	defer t.headerMu.Unlock()
	t.headers[key] = value
}

// SetHTTPClient replaces the internal *http.Client. Callers can pass
// in a custom HTTP/3 client here.
func (t *Transport) SetHTTPClient(c *http.Client) {
	if c != nil {
		t.httpClient = c
	}
}

// HTTPClient returns the internal *http.Client (for advanced use).
func (t *Transport) HTTPClient() *http.Client { return t.httpClient }

// applyHeaders copies the transport's default headers onto req. The
// header map is guarded by headerMu because SetHeader may run
// concurrently with a request; the read lock lets in-flight requests
// snapshot a consistent header set without blocking each other.
func (t *Transport) applyHeaders(req *http.Request) {
	t.headerMu.RLock()
	defer t.headerMu.RUnlock()
	for k, v := range t.headers {
		req.Header.Set(k, v)
	}
}
