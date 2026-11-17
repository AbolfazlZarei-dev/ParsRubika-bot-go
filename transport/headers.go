package transport

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/proxy"
)

// ErrProxyUnavailable is returned by a fail-closed transport when a
// proxy was requested but could not be configured (bad URL, unsupported
// scheme, dialer error). Rather than silently falling back to a direct
// connection — which would bypass the user's proxy and could leak the
// auth key to an unintended host — the transport's DialContext returns
// this error for every request until the proxy is fixed.
var ErrProxyUnavailable = errors.New("transport: proxy requested but unavailable; refusing to fall back to a direct connection")

// DefaultHeaders returns the default HTTP header set used by ParsRubika.
func DefaultHeaders(userAgent string) map[string]string {
	return map[string]string{
		"origin":       "https://m.rubika.ir",
		"referer":      "https://m.rubika.ir/",
		"content-type": "application/json",
		"connection":   "keep-alive",
		"user-agent":   userAgent,
	}
}

// AndroidHeaders returns the header set used when the platform is
// Android. The User-Agent is overridden to `okhttp/3.12.1` and the
// Origin/Referer headers are dropped.
func AndroidHeaders(userAgent string) map[string]string {
	return map[string]string{
		"content-type": "application/json",
		"connection":   "keep-alive",
		"user-agent":   "okhttp/3.12.1",
	}
}

// NewHTTPClient builds an *http.Client tuned for high-throughput API
// traffic against a small number of hosts:
//
//   - HTTP/2 is force-attempted over TLS (multiplexes many requests over
//     one connection, removing per-request handshake cost).
//   - A generous idle-connection pool (256 per host, 512 overall) keeps
//     warm connections ready under bursty concurrency, so requests skip
//     the TCP+TLS handshake on the hot path.
//   - Tight handshake/header deadlines fail fast on stalled peers instead
//     of burning the whole request timeout.
//   - Optional proxy support — accepts HTTP, HTTPS, and SOCKS5 proxy
//     URLs (e.g. "http://127.0.0.1:8080", "socks5://127.0.0.1:1080",
//     "socks5://user:pass@127.0.0.1:1080").
//
// The returned client is safe for concurrent use.
func NewHTTPClient(timeout time.Duration, proxyURL string, insecure bool) *http.Client {
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	tr := &http.Transport{
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: insecure},
		DialContext:         dialer.DialContext,
		MaxIdleConns:        512,
		MaxIdleConnsPerHost: 256,
		// Bound the number of LIVE (in-flight + idle) connections per host
		// so a misbehaving server cannot keep an unbounded number of
		// sockets open. 64 is ample for the Rubika API host under heavy
		// concurrency; excess requests queue instead of opening new conns.
		MaxConnsPerHost:       64,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ResponseHeaderTimeout: 30 * time.Second,
		ForceAttemptHTTP2:     true,
	}
	if proxyURL != "" {
		if err := applyProxy(tr, proxyURL); err != nil {
			// Fail CLOSED: never fall back to a direct connection when a
			// proxy was explicitly requested. Install a dialer that
			// refuses every connect so the misconfiguration surfaces as
			// request errors instead of silently leaking traffic.
			tr.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
				return nil, ErrProxyUnavailable
			}
		}
	}
	return &http.Client{
		Transport: tr,
		Timeout:   timeout,
		// Rubika API endpoints never redirect. A redirect on the auth path
		// would re-send the JSON body (which carries the auth key) to a
		// different host; net/http only strips the Authorization header on
		// redirect, not a JSON body. Disabling redirects entirely closes
		// that leakage path. (User-URL fetches use netutil.NewSafeClient,
		// which still follows redirects through the SSRF-controlled dialer.)
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

// applyProxy configures the supplied http.Transport to use the given
// proxy URL. It supports HTTP, HTTPS, and SOCKS5 proxies (with or
// without authentication). It returns an error if the proxy URL is
// malformed or the SOCKS5 dialer cannot be built; the caller is
// expected to fail closed in that case.
//
// Supported formats:
//
//	http://[user:pass@]host:port
//	https://[user:pass@]host:port
//	socks5://[user:pass@]host:port
//	socks5h://[user:pass@]host:port (DNS via proxy)
//
// For SOCKS5 we install a custom DialContext so the rest of the
// transport (TLS, HTTP/2, connection pooling) works unchanged.
func applyProxy(tr *http.Transport, proxyURL string) error {
	u, err := url.Parse(proxyURL)
	if err != nil {
		return fmt.Errorf("transport: parse proxy url: %w", err)
	}
	scheme := strings.ToLower(u.Scheme)
	switch scheme {
	case "socks5", "socks5h":
		if u.Host == "" {
			return fmt.Errorf("transport: socks5 proxy missing host")
		}
		// Build a SOCKS5 dialer with optional auth.
		var auth *proxy.Auth
		if u.User != nil {
			pw, _ := u.User.Password()
			auth = &proxy.Auth{
				User:     u.User.Username(),
				Password: pw,
			}
		}
		dialer, err := proxy.SOCKS5("tcp", u.Host, auth, &net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		})
		if err != nil {
			return fmt.Errorf("transport: socks5 dialer: %w", err)
		}
		// Wrap as DialContext so http.Transport can use it.
		tr.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			// SOCKS5 dialers don't honour context natively; we just
			// call Dial and rely on the dialer's own Timeout.
			return dialer.Dial(network, addr)
		}
		// When scheme is socks5h, DNS is resolved by the proxy
		// (already the default for proxy.SOCKS5). For plain socks5,
		// Go still resolves locally first — both work in practice.
	default:
		// HTTP / HTTPS proxy.
		tr.Proxy = http.ProxyURL(u)
	}
	return nil
}
