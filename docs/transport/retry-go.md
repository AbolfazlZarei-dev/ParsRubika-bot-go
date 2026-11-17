# `transport` — Retry

_Source: [`retry.go`](../../ParsRubika/transport/retry.go)_

## Symbols

### `RetryPolicy`

```go
type RetryPolicy struct
```

RetryPolicy describes how a failed request should be retried.

### `DefaultRetryPolicy`

```go
func DefaultRetryPolicy() RetryPolicy
```

DefaultRetryPolicy returns the retry policy used by ParsRubika by
default for idempotent requests (GET): 5 attempts, 1-second base
delay, capped at 30 seconds, and retry on transient network errors
and 5xx responses. Because GET is idempotent, retrying on an
ambiguous transport error (the server may or may not have seen the
request) is safe.

### `PostRetryPolicy`

```go
func PostRetryPolicy() RetryPolicy
```

PostRetryPolicy is the retry policy for non-idempotent POST requests
(every Rubika message-send, file-upload-start, etc.). It retries ONLY
on explicit, server-signalled retryable HTTP statuses (408, 425, 429,
500, 502, 503, 504) and NEVER on a transport-level error
(status 0 / err != nil).

A transport error on a POST is ambiguous: the request may have
reached the server and been applied (e.g. the message may already be
delivered) even though the client never saw the response. Retrying
in that case risks delivering a DUPLICATE message, so the request is
failed back to the caller instead. The server-signalled statuses are
safe to retry because the server explicitly states it did not apply
the request (or is asking the client to slow down).

### `Sleep`

```go
func (p RetryPolicy) Sleep(ctx context.Context, attempt int) error
```

Sleep waits for the exponential backoff of the supplied attempt
number, with up to ±25 % jitter so that a burst of retries from many
clients does not synchronise (thundering herd). The first attempt
uses base delay; subsequent attempts double the delay each time,
capped by `MaxDelay`.

### `ShouldRetry`

```go
func (p RetryPolicy) ShouldRetry(status int, err error) bool
```

ShouldRetry reports whether the policy says the request should be
retried for the given (status, err) pair.

### `ErrMaxRetries`

```go
var ErrMaxRetries = errors.New("transport: max retries exceeded")
```

ErrMaxRetries is returned when a request exhausts all retries.
