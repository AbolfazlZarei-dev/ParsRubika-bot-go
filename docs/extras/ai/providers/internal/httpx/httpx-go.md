# `httpx` — Package httpx is a tiny shared HTTP shim for the provider packages

_Source: [`httpx.go`](../../../../../../ParsRubika/extras/ai/providers/internal/httpx/httpx.go)_

## Symbols

### `MaxBodyBytes`

```go
const MaxBodyBytes = 8 << 20
```

MaxBodyBytes caps a response body so a misbehaving endpoint cannot
exhaust memory.

### `PostJSON`

```go
func PostJSON(ctx context.Context, client *http.Client, url string, headers map[string]string, body any) (int, []byte, error)
```

PostJSON POSTs the marshalled body to url with the given headers and
returns the HTTP status code and the raw response body. A nil client
falls back to http.DefaultClient. The Content-Type header defaults to
application/json when not set in headers.

### `Truncate`

```go
func Truncate(b []byte, n int) string
```

Truncate returns at most n bytes of b as a string, with an ellipsis
when it was shortened. Used to keep error messages readable.
