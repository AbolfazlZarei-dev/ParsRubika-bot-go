# `account` — Lifecycle

_Source: [`lifecycle.go`](../../ParsRubika/account/lifecycle.go)_

## Symbols

### `Start`

```go
func (a *Account) Start(ctx context.Context, phone string) error
```

Start performs the connection + (optional) sign-in handshake. If a
session has been saved previously it is reused; otherwise the
`phoneNumber` (or the one passed in) is used to send a code and
sign in interactively.

### `Stop`

```go
func (a *Account) Stop()
```

Stop tears down the WebSocket and the transport. It is idempotent and
signals the reconnect loop not to redial afterwards.

### `Run`

```go
func (a *Account) Run(ctx context.Context, phone string) error
```

Run is the convenience method that calls Start + Listen. It blocks
until `ctx` is cancelled or Stop/Close is called.

### `Listen`

```go
func (a *Account) Listen(ctx context.Context) error
```

Listen starts the WebSocket update loop. It blocks until `ctx` is
cancelled or a fatal error occurs.

Reconnect strategy: exponential backoff with a cap and full jitter,
reset to zero after a successful handshake. The previous fixed 3 s
sleep could hammer the server (and burn CPU) during a persistent
outage; the new strategy slows down (up to wsMaxBackoff) and recovers
immediately once the connection is healthy again.
