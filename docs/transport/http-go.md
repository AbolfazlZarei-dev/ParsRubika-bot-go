# `transport` — Http

_Source: [`http.go`](../../ParsRubika/transport/http.go)_

## Symbols

### `MaxResponseBody`

```go
const MaxResponseBody = 16 << 20 // 16 MiB
```

MaxResponseBody caps how many bytes of an HTTP response body the
transport will read into memory. Rubika API responses are small JSON
documents; the cap exists so a malicious or buggy server cannot make
the client allocate unbounded memory by streaming an enormous body.
A response that exceeds the cap fails with ErrResponseBodyTooLarge.

### `ErrResponseBodyTooLarge`

```go
var ErrResponseBodyTooLarge = fmt.Errorf("transport: response body exceeds %d bytes", MaxResponseBody)
```

ErrResponseBodyTooLarge is returned when a response body exceeds
MaxResponseBody.

### `Post`

```go
func (t *Transport) Post(ctx context.Context, url string, payload any) ([]byte, int, error)
```

Post sends a POST request with the supplied JSON payload and returns
the raw response body. It retries on transient failures according to
the configured RetryPolicy.

The payload is encoded into a process-wide pooled buffer rather than a
freshly allocated slice, so the encode step reuses warm memory across
requests. The resulting bytes are handed straight to bytes.NewReader
in postRaw with no intermediate copy.

### `Get`

```go
func (t *Transport) Get(ctx context.Context, url string) ([]byte, int, error)
```

Get sends a GET request and returns the raw response body.

### `PostJSON`

```go
func (t *Transport) PostJSON(ctx context.Context, url string, payload any, out any) (int, error)
```

PostJSON is a convenience helper that JSON-decodes the response body
into `out`.

### `New`

```go
func New(opts ...Option) *Transport
```

New constructs a Transport with the supplied options. It uses the
default HTTP/2-enabled client and the standard header set.

### `SetHeader`

```go
func (t *Transport) SetHeader(key, value string)
```

SetHeader updates a single header. Useful for the bot API client
which uses a different base URL per token. It is safe to call
concurrently with in-flight requests.

### `SetHTTPClient`

```go
func (t *Transport) SetHTTPClient(c *http.Client)
```

SetHTTPClient replaces the internal *http.Client. Callers can pass
in a custom HTTP/3 client here.

### `HTTPClient`

```go
func (t *Transport) HTTPClient() *http.Client { return t.httpClient }
```

HTTPClient returns the internal *http.Client (for advanced use).
