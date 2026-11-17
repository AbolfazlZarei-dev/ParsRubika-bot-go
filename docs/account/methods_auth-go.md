# `account` — Methods auth

_Source: [`methods_auth.go`](../../ParsRubika/account/methods_auth.go)_

## Symbols

### `SendCode`

```go
func (a *Account) SendCode(ctx context.Context, phone, passKey, sendType string) (*events.Event, error)
```

SendCode triggers the SMS / Internal verification code for `phone`.

### `SignIn`

```go
func (a *Account) SignIn(ctx context.Context, phoneCode, phone, phoneCodeHash, publicKey string) (*events.Event, error)
```

SignIn completes the registration handshake.

### `RegisterDevice`

```go
func (a *Account) RegisterDevice(ctx context.Context, deviceModel string) (*events.Event, error)
```

RegisterDevice registers the device model with the server.

### `UnregisterDevice`

```go
func (a *Account) UnregisterDevice(ctx context.Context) (*events.Event, error)
```

UnregisterDevice removes the current device registration from the
server. Use this for a clean logout (without it the device lingers
in the server's device list and may keep receiving push tokens).
For a full logout, follow this with Logout (in methods_auth2fa.go).

### `Logout`

```go
func (a *Account) Logout(ctx context.Context) (*events.Event, error)
```

Logout logs out the current session.

### `DeleteAccount`

```go
func (a *Account) DeleteAccount(ctx context.Context) (*events.Event, error)
```

DeleteAccount requests account deletion.

### `UpgradeToAPI6`

```go
func (a *Account) UpgradeToAPI6(ctx context.Context) (privateKeyPEM string, err error)
```

UpgradeToAPI6 upgrades an API5 auth to API6 by uploading a new RSA
public key. Returns the private key (PEM) that the caller should
store alongside the auth.

### `GetUnconfirmedSessions`

```go
func (a *Account) GetUnconfirmedSessions(ctx context.Context) (*events.Event, error)
```

GetUnconfirmedSessions returns the unconfirmed session list.

### `ActionOnUnconfirmedSession`

```go
func (a *Account) ActionOnUnconfirmedSession(ctx context.Context, action, unconfirmedSessionKey string) (*events.Event, error)
```

ActionOnUnconfirmedSession confirms or denies an unconfirmed session.

### `ConfirmUnconfirmedSession`

```go
func (a *Account) ConfirmUnconfirmedSession(ctx context.Context, key string) (*events.Event, error)
```

ConfirmUnconfirmedSession is a convenience wrapper (action = "Confirm").

### `DenyUnconfirmedSession`

```go
func (a *Account) DenyUnconfirmedSession(ctx context.Context, key string) (*events.Event, error)
```

DenyUnconfirmedSession is a convenience wrapper (action = "Deny").

### `TerminateSession`

```go
func (a *Account) TerminateSession(ctx context.Context, sessionID string) (*events.Event, error)
```

TerminateSession terminates a single active session by id.

### `TerminateOtherSessions`

```go
func (a *Account) TerminateOtherSessions(ctx context.Context) (*events.Event, error)
```

TerminateOtherSessions terminates every other active session.

### `GetMySessions`

```go
func (a *Account) GetMySessions(ctx context.Context) (*events.Event, error)
```

GetMySessions returns the active sessions list.
