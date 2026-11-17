# `account` — Request

_Source: [`request.go`](../../ParsRubika/account/request.go)_

## Symbols

### `RequestOptions`

```go
type RequestOptions struct
```

RequestOptions carries the optional flags for the encrypted request.

### `Request`

```go
func (a *Account) Request(ctx context.Context, opts RequestOptions) (*events.Event, error)
```

Request is the central encrypted request dispatcher. It mirrors
the original Builder: it encrypts the input + sends it through the
transport layer, decrypts the response, validates `status_det ==
OK` and returns a typed *events.Event.
