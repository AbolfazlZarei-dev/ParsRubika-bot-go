// Package transport provides the network layer for ParsRubika. It
// exposes an HTTP client with built-in HTTP/2 (and optional HTTP/3),
// connection pooling, retries with exponential backoff, rate limiting,
// WebSocket support and the Rubika data-center discovery protocol.
//
// The package is split across several small files so each concern can
// be inspected and tested independently:
//
// - transport.go : the Transport struct + configuration
// - http2.go : the HTTP/2-enabled http.Client
// - (HTTP/3 lives in the optional http3 module)
// - websocket.go : the WebSocket connection used for live updates
// - dcs.go : data-center discovery
// - retry.go : retry policy + exponential backoff
// - ratelimit.go : per-token rate limiter
// - headers.go : the default request headers
package transport

import (
	"net/http"
	"sync"
	"time"
)

// Transport holds the shared HTTP client and connection-pool
// configuration. It is safe for concurrent use after construction.
type Transport struct {
	httpClient *http.Client
	headers    map[string]string
	// headerMu guards the headers map. SetHeader may be called
	// concurrently with in-flight requests (which call applyHeaders),
	// so both the writer and the readers take the lock.
	headerMu   sync.RWMutex
	userAgent  string
	proxy      string
	timeout    time.Duration
	maxRetries int
}

// Option configures a Transport at construction time.
type Option func(*Transport)

// WithUserAgent overrides the default User-Agent string.
func WithUserAgent(ua string) Option {
	return func(t *Transport) { t.userAgent = ua }
}

// WithProxy sets an HTTP(S) proxy URL.
func WithProxy(proxy string) Option {
	return func(t *Transport) { t.proxy = proxy }
}

// WithTimeout sets the HTTP timeout.
func WithTimeout(d time.Duration) Option {
	return func(t *Transport) { t.timeout = d }
}

// WithMaxRetries overrides the default max-retries value.
func WithMaxRetries(n int) Option {
	return func(t *Transport) { t.maxRetries = n }
}
