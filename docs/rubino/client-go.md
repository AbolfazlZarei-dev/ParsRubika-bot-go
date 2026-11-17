# `rubino` — Everything-in-one client (App) and one-line helpers

_Source: [`client.go`](../../ParsRubika/rubino/client.go)_

## Symbols

### `RubinoAPIURL`

```go
const RubinoAPIURL = "https://rubino1.iranlms.ir/"
```

RubinoAPIURL is the base URL for the Rubino API.

### `Rubino`

```go
type Rubino struct
```

Rubino is the Rubino social-network client.

### `Option`

```go
type Option func(*Rubino)
```

Option configures a Rubino at construction time.

### `WithAuth`

```go
func WithAuth(auth string) Option
```

WithAuth sets the auth key.

### `WithUserAgent`

```go
func WithUserAgent(ua string) Option
```

WithUserAgent overrides the User-Agent string.

### `WithProxy`

```go
func WithProxy(proxy string) Option
```

WithProxy sets a proxy URL.

### `WithTimeout`

```go
func WithTimeout(seconds int) Option
```

WithTimeout sets the HTTP timeout in seconds.

### `WithLogger`

```go
func WithLogger(l *log.Logger) Option
```

WithLogger injects a custom logger.

### `WithProfile`

```go
func WithProfile(profileID string) Option
```

WithProfile sets the active profile ID. Some methods use this
default when their `profileID` parameter is empty.

### `New`

```go
func New(opts ...Option) *Rubino
```

New constructs a Rubino client.

### `Close`

```go
func (s *Rubino) Close() error { return nil }
```

Close releases the underlying transport.

### `SetHTTPClient`

```go
func (s *Rubino) SetHTTPClient(c *http.Client)
```

SetHTTPClient lets callers swap in an HTTP/3 client.

### `SetAuth`

```go
func (s *Rubino) SetAuth(auth string)
```

SetAuth updates the auth key. Used by callers that want to reuse
the account's auth key (after the account signs in) instead of a
separate Rubino login — Rubino accepts the account's auth key.
