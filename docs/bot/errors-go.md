# `bot` — Errors

_Source: [`errors.go`](../../ParsRubika/bot/errors.go)_

## Symbols

### `APIError`

```go
type APIError struct
```

APIError is returned when the Bot API answers with a non-OK status.

### `Error`

```go
func (e *APIError) Error() string
```

Error implements the error interface.

### `Temporary`

```go
func (e *APIError) Temporary() bool
```

Temporary reports whether retrying the call might succeed.

### `IsStatus`

```go
func IsStatus(err error, status string) bool
```

IsStatus reports whether err (or any error it wraps) is an *APIError
carrying the given status string.

### `IsRateLimited`

```go
func IsRateLimited(err error) bool
```

IsRateLimited reports whether err is a "too many requests" failure.

### `NewAPIError`

```go
func NewAPIError(status, devMessage string) *APIError
```

NewAPIError is a convenience constructor for callers that want to
synthesise an error from a custom status string.

### `ErrContinue`

```go
var ErrContinue = errors.New("bot: continue to next handler")
```

ErrContinue can be returned by a handler to let the dispatcher keep
looking for the next matching handler instead of stopping at this one.

### `ErrStopped`

```go
var ErrStopped = errors.New("bot: stopped")
```

ErrStopped is returned by operations attempted after the bot stopped.

### `ErrorContext`

```go
type ErrorContext struct
```

ErrorContext describes a failure raised while processing an update.

### `Error`

```go
func (e *ErrorContext) Error() string
```

Error summarises the failure.

### `Printf`

```go
func (s stdLogger) Printf(format string, args ...any) { s.l.Printf(format, args...) }
```
