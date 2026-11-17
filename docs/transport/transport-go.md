# `transport` — Package transport provides the network layer for ParsRubika. It

_Source: [`transport.go`](../../ParsRubika/transport/transport.go)_

## Symbols

### `Transport`

```go
type Transport struct
```

Transport holds the shared HTTP client and connection-pool
configuration. It is safe for concurrent use after construction.

### `Option`

```go
type Option func(*Transport)
```

Option configures a Transport at construction time.

### `WithUserAgent`

```go
func WithUserAgent(ua string) Option
```

WithUserAgent overrides the default User-Agent string.

### `WithProxy`

```go
func WithProxy(proxy string) Option
```

WithProxy sets an HTTP(S) proxy URL.

### `WithTimeout`

```go
func WithTimeout(d time.Duration) Option
```

WithTimeout sets the HTTP timeout.

### `WithMaxRetries`

```go
func WithMaxRetries(n int) Option
```

WithMaxRetries overrides the default max-retries value.
