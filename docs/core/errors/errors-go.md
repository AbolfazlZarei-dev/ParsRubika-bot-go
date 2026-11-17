# `errors` — Package errors defines the typed and sentinel errors used across

_Source: [`errors.go`](../../../ParsRubika/core/errors/errors.go)_

## Symbols

### `Kind`

```go
type Kind int
```

Kind classifies an error so that callers can react without parsing
the message text.

### `KindUnknown`

```go
const KindUnknown Kind = iota
```

### `KindNetwork`

```go
const KindNetwork
```

### `KindTimeout`

```go
const KindTimeout
```

### `KindAuth`

```go
const KindAuth
```

### `KindValidation`

```go
const KindValidation
```

### `KindAPI`

```go
const KindAPI
```

### `KindUpload`

```go
const KindUpload
```

### `KindDownload`

```go
const KindDownload
```

### `KindSession`

```go
const KindSession
```

### `KindConfig`

```go
const KindConfig
```

### `KindNotImplemented`

```go
const KindNotImplemented
```

### `KindCancelled`

```go
const KindCancelled
```

### `KindStopDispatch`

```go
const KindStopDispatch
```

### `TypedError`

```go
type TypedError struct
```

TypedError is the structured error returned by ParsRubika. It carries
a Kind, a human-readable message, and an optional server-side status
string (e.g. "InvalidAuth", "NotRegistered").

### `Error`

```go
func (e *TypedError) Error() string
```

Error implements the error interface.

### `Unwrap`

```go
func (e *TypedError) Unwrap() error { return e.Wrapped }
```

Unwrap supports errors.Is / errors.As.

### `String`

```go
func (k Kind) String() string
```

String implements fmt.Stringer for Kind.

### `New`

```go
func New(kind Kind, msg string) *TypedError
```

New constructs a fresh TypedError.

### `Wrap`

```go
func Wrap(kind Kind, msg string, inner error) *TypedError
```

Wrap returns a TypedError that wraps `inner` while preserving its
kind for `errors.Is` matching.

### `WithStatus`

```go
func WithStatus(kind Kind, status, msg string) *TypedError
```

WithStatus attaches a server-side status string to the error.

### `ErrNetwork`

```go
var ErrNetwork = New(KindNetwork, "network error")
```

Common sentinel errors for use with `errors.Is`.

### `ErrTimeout`

```go
var ErrTimeout = New(KindTimeout, "request timed out")
```

Common sentinel errors for use with `errors.Is`.

### `ErrInvalidAuth`

```go
var ErrInvalidAuth = WithStatus(KindAuth, "InvalidAuth", "invalid auth key")
```

Common sentinel errors for use with `errors.Is`.

### `ErrNotRegistered`

```go
var ErrNotRegistered = WithStatus(KindAuth, "NotRegistered", "user not registered")
```

Common sentinel errors for use with `errors.Is`.

### `ErrCodeUsed`

```go
var ErrCodeUsed = WithStatus(KindAuth, "CodeIsUsed", "verification code already used")
```

Common sentinel errors for use with `errors.Is`.

### `ErrCodeExpired`

```go
var ErrCodeExpired = WithStatus(KindAuth, "CodeIsExpired", "verification code expired")
```

Common sentinel errors for use with `errors.Is`.

### `ErrTooManyRequests`

```go
var ErrTooManyRequests = WithStatus(KindAPI, "TooRequests", "too many requests")
```

Common sentinel errors for use with `errors.Is`.

### `ErrServerError`

```go
var ErrServerError = WithStatus(KindAPI, "ServerError", "server error")
```

Common sentinel errors for use with `errors.Is`.

### `ErrURLNotFound`

```go
var ErrURLNotFound = WithStatus(KindAPI, "UrlNotFound", "url not found")
```

Common sentinel errors for use with `errors.Is`.

### `ErrInvalidInput`

```go
var ErrInvalidInput = WithStatus(KindValidation, "InvalidInput", "invalid input")
```

Common sentinel errors for use with `errors.Is`.

### `ErrNoConnection`

```go
var ErrNoConnection = WithStatus(KindNetwork, "NoConnection", "no connection")
```

Common sentinel errors for use with `errors.Is`.

### `ErrUnsupportedAPI`

```go
var ErrUnsupportedAPI = WithStatus(KindAPI, "NotSupportedApiVersion", "unsupported API version")
```

Common sentinel errors for use with `errors.Is`.

### `ErrUploadFailed`

```go
var ErrUploadFailed = New(KindUpload, "upload failed")
```

Common sentinel errors for use with `errors.Is`.

### `ErrDownloadFailed`

```go
var ErrDownloadFailed = New(KindDownload, "download failed")
```

Common sentinel errors for use with `errors.Is`.

### `ErrSessionEmpty`

```go
var ErrSessionEmpty = New(KindSession, "session is empty")
```

Common sentinel errors for use with `errors.Is`.

### `ErrSessionCorrupt`

```go
var ErrSessionCorrupt = New(KindSession, "session data is corrupt")
```

Common sentinel errors for use with `errors.Is`.

### `ErrConfigInvalid`

```go
var ErrConfigInvalid = New(KindConfig, "configuration is invalid")
```

Common sentinel errors for use with `errors.Is`.

### `ErrCancelled`

```go
var ErrCancelled = New(KindCancelled, "operation cancelled")
```

Common sentinel errors for use with `errors.Is`.

### `ErrStopDispatch`

```go
var ErrStopDispatch = New(KindStopDispatch, "stop dispatching")
```

Common sentinel errors for use with `errors.Is`.

### `ErrNotImplemented`

```go
var ErrNotImplemented = New(KindNotImplemented, "not implemented")
```

Common sentinel errors for use with `errors.Is`.

### `FromStatus`

```go
func FromStatus(status, msg string) error
```

FromStatus maps a server-side `status_det` string to the matching
sentinel error. Unknown statuses produce a generic *TypedError with
the original status preserved.

### `IsKind`

```go
func IsKind(err error, k Kind) bool
```

IsKind reports whether err is a *TypedError of the given Kind.

### `IsStatus`

```go
func IsStatus(err error, status string) bool
```

IsStatus reports whether err is a *TypedError with the given status.
