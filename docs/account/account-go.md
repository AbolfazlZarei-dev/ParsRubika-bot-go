# `account` — Client, options and constructor

_Source: [`account.go`](../../ParsRubika/account/account.go)_

## Symbols

### `Account`

```go
type Account struct
```

Account is the high-level user/account client.

### `Listener`

```go
type Listener func(ctx context.Context, e *events.Event) error
```

Listener is the signature of a registered update callback.

### `CodeProvider`

```go
type CodeProvider func(ctx context.Context, hint string) (string, error)
```

CodeProvider is the signature of the callback used by the interactive
sign-in flow to obtain the SMS code (or 2FA password) from the user.
The `hint` argument is a short human-readable description of what is
being requested (e.g. "Code:" or "Password [hint] > "). The provider
should return the user's response, trimmed of trailing whitespace.

When no CodeProvider is set, sign-in falls back to reading from stdin
(interactive mode). Setting a CodeProvider makes the account usable
in headless environments (servers, Docker, systemd).

### `Matcher`

```go
type Matcher interface
```

Matcher is the contract implemented by every account-side matcher
(text, command, chat-type, ...). The matchers live in the
`parsrubika/account/matchers` sub-package.

### `Option`

```go
type Option func(*Account) error
```

Option configures an Account at construction time.

### `WithAuth`

```go
func WithAuth(auth string) Option
```

WithAuth sets the auth key.

### `WithPrivateKey`

```go
func WithPrivateKey(pk string) Option
```

WithPrivateKey sets the RSA private key (PEM-encoded).

### `WithPhoneNumber`

```go
func WithPhoneNumber(phone string) Option
```

WithPhoneNumber sets the phone number used during the registration
handshake.

### `WithUserAgent`

```go
func WithUserAgent(ua string) Option
```

WithUserAgent overrides the default User-Agent string.

### `WithTimeout`

```go
func WithTimeout(d time.Duration) Option
```

WithTimeout sets the HTTP timeout.

### `WithLangCode`

```go
func WithLangCode(code string) Option
```

WithLangCode sets the language code (default "fa").

### `WithParseMode`

```go
func WithParseMode(mode enums.ParseMode) Option
```

WithParseMode sets the default parse mode.

### `WithProxy`

```go
func WithProxy(proxy string) Option
```

WithProxy sets an HTTP(S) proxy URL.

### `WithPlatform`

```go
func WithPlatform(p string) Option
```

WithPlatform selects "PWA" (default) or "Android".

### `WithMaxRetries`

```go
func WithMaxRetries(n int) Option
```

WithMaxRetries overrides the default max-retries value.

### `WithSequentialListeners`

```go
func WithSequentialListeners(b bool) Option
```

WithSequentialListeners, when true, makes the dispatcher run
listeners sequentially and stop after the first match.

### `WithLogger`

```go
func WithLogger(l core.Logger) Option
```

WithLogger injects a custom logger.

### `WithDebug`

```go
func WithDebug(b bool) Option
```

WithDebug enables verbose request/response/WebSocket logging for
troubleshooting. Off by default — when disabled, plaintext payloads
are never logged.

### `WithCodeProvider`

```go
func WithCodeProvider(p CodeProvider) Option
```

WithCodeProvider installs a callback used by the interactive sign-in
flow to obtain the SMS code (or 2FA password) from the caller. When
set, the account client can run headlessly (servers, Docker, systemd)
instead of prompting on stdin/stdout.

Example:

acc, _ := account.New(
account.WithPhoneNumber("+989123456789"),
account.WithCodeProvider(func(ctx context.Context, hint string) (string, error) {
// hint is "Code:" or "Password [hint] > "
return myHTTPUI.PromptUser(ctx, hint)
}),
)

### `WithStore`

```go
func WithStore(s storage.Store) Option
```

WithStore lets the caller supply a custom Store implementation.

### `WithSQLiteSession`

```go
func WithSQLiteSession(name string) Option
```

WithSQLiteSession uses a SQLite session file with the given name.

The session is stored in plain text. For at-rest encryption of the
auth key and RSA private key use WithSQLiteSessionPassphrase instead.

### `WithSQLiteSessionPassphrase`

```go
func WithSQLiteSessionPassphrase(name, passphrase string) Option
```

WithSQLiteSessionPassphrase uses a SQLite session file with the given
name and enables AES-256-GCM at-rest encryption for the auth and
private_key columns. The passphrase is NOT stored alongside the
session file; losing it makes the session unrecoverable.

### `WithFileSession`

```go
func WithFileSession(name string) Option
```

WithFileSession uses a JSON file session with the given name.

The session is stored in plain text. For at-rest encryption of the
auth key and RSA private key use WithFileSessionPassphrase instead.

### `WithFileSessionPassphrase`

```go
func WithFileSessionPassphrase(name, passphrase string) Option
```

WithFileSessionPassphrase uses a JSON file session with the given name
and enables AES-256-GCM at-rest encryption for the auth and
private_key fields. The passphrase is NOT stored alongside the session
file; losing it makes the session unrecoverable.

### `WithHTTPClient`

```go
func WithHTTPClient(c *http.Client) Option
```

WithHTTPClient lets the caller supply a pre-configured *http.Client
(e.g. a HTTP/3 client). When set, the supplied client is used as-is
and the proxy/timeout options are ignored.

### `New`

```go
func New(opts ...Option) (*Account, error)
```

New constructs an Account using the supplied options.

### `Name`

```go
func (a *Account) Name() string { return a.name }
```

Name returns the session name.

### `Auth`

```go
func (a *Account) Auth() string { return a.auth }
```

Auth returns the configured auth key.

### `GUID`

```go
func (a *Account) GUID() string { return a.guid }
```

GUID returns the GUID of the logged-in user.

### `ParseMode`

```go
func (a *Account) ParseMode() enums.ParseMode { return a.parseMode }
```

ParseMode returns the default parse mode.

### `UserAgent`

```go
func (a *Account) UserAgent() string { return a.userAgent }
```

UserAgent returns the configured User-Agent string.

### `Store`

```go
func (a *Account) Store() storage.Store { return a.store }
```

Store returns the session store.

### `Formatter`

```go
func (a *Account) Formatter() text.Formatter { return a.formatter }
```

Formatter returns the text formatter.

### `Uploader`

```go
func (a *Account) Uploader() *files.Uploader { return a.uploader }
```

Uploader returns the file uploader.

### `Downloader`

```go
func (a *Account) Downloader() *files.Downloader { return a.downloader }
```

Downloader returns the file downloader.

### `Close`

```go
func (a *Account) Close() error
```

Close releases the WebSocket + transport + store. It is idempotent and
signals the reconnect loop (if running) not to redial afterwards.
