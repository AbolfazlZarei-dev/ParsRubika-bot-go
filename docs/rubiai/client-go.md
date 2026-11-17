# `rubiai` — Everything-in-one client (App) and one-line helpers

_Source: [`client.go`](../../ParsRubika/rubiai/client.go)_

## Symbols

### `BaseURL`

```go
const BaseURL = "https://rubichat1.iranlms.ir/"
```

BaseURL is the default RubiAI API endpoint.

### `ClientInfo`

```go
var ClientInfo = map[string]any{
```

ClientInfo is the `client` block sent inside every RubiAI request. It
identifies the caller as the official Android Rubika app so the server
accepts the request.

### `RubiaAI`

```go
type RubiaAI struct
```

RubiaAI is the RubiAI assistant client.

### `Option`

```go
type Option func(*RubiaAI)
```

Option configures a RubiaAI client at construction time.

### `WithAuth`

```go
func WithAuth(auth string) Option
```

WithAuth sets the RubiAI auth key. Required for every request.

### `WithUserAgent`

```go
func WithUserAgent(ua string) Option
```

WithUserAgent overrides the default User-Agent string.

### `WithProxy`

```go
func WithProxy(proxy string) Option
```

WithProxy sets a proxy URL (http://, https://, socks5://, socks5h://).

### `WithTimeout`

```go
func WithTimeout(d time.Duration) Option
```

WithTimeout sets the HTTP timeout.

### `WithMaxRetries`

```go
func WithMaxRetries(n int) Option
```

WithMaxRetries overrides the default retry count for transient errors.

### `WithLogger`

```go
func WithLogger(l *log.Logger) Option
```

WithLogger injects a custom logger.

### `New`

```go
func New(opts ...Option) *RubiaAI
```

New constructs a RubiaAI client.

### `SetHTTPClient`

```go
func (r *RubiaAI) SetHTTPClient(c *http.Client)
```

SetHTTPClient lets callers swap in a pre-configured *http.Client
(e.g. an HTTP/3 client or one bound to a specific proxy).

### `SetAuth`

```go
func (r *RubiaAI) SetAuth(auth string)
```

SetAuth updates the auth key. Useful when reusing the account client's
auth key after the account signs in (RubiAI accepts the account auth).

### `Close`

```go
func (r *RubiaAI) Close() error { return nil }
```

Close releases the underlying transport.
