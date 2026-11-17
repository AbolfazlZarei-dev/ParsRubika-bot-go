# `transport` — Headers

_Source: [`headers.go`](../../ParsRubika/transport/headers.go)_

## Symbols

### `ErrProxyUnavailable`

```go
var ErrProxyUnavailable = errors.New("transport: proxy requested but unavailable; refusing to fall back to a direct connection")
```

ErrProxyUnavailable is returned by a fail-closed transport when a
proxy was requested but could not be configured (bad URL, unsupported
scheme, dialer error). Rather than silently falling back to a direct
connection — which would bypass the user's proxy and could leak the
auth key to an unintended host — the transport's DialContext returns
this error for every request until the proxy is fixed.

### `DefaultHeaders`

```go
func DefaultHeaders(userAgent string) map[string]string
```

DefaultHeaders returns the default HTTP header set used by ParsRubika.

### `AndroidHeaders`

```go
func AndroidHeaders(userAgent string) map[string]string
```

AndroidHeaders returns the header set used when the platform is
Android. The User-Agent is overridden to `okhttp/3.12.1` and the
Origin/Referer headers are dropped.

### `NewHTTPClient`

```go
func NewHTTPClient(timeout time.Duration, proxyURL string, insecure bool) *http.Client
```

NewHTTPClient builds an *http.Client tuned for high-throughput API
traffic against a small number of hosts:

- HTTP/2 is force-attempted over TLS (multiplexes many requests over
one connection, removing per-request handshake cost).
- A generous idle-connection pool (256 per host, 512 overall) keeps
warm connections ready under bursty concurrency, so requests skip
the TCP+TLS handshake on the hot path.
- Tight handshake/header deadlines fail fast on stalled peers instead
of burning the whole request timeout.
- Optional proxy support — accepts HTTP, HTTPS, and SOCKS5 proxy
URLs (e.g. "http://127.0.0.1:8080", "socks5://127.0.0.1:1080",
"socks5://user:pass@127.0.0.1:1080").

The returned client is safe for concurrent use.
