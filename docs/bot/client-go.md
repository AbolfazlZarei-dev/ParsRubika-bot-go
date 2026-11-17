# `bot` — Everything-in-one client (App) and one-line helpers

_Source: [`client.go`](../../ParsRubika/bot/client.go)_

## Symbols

### `Middleware`

```go
type Middleware func(ctx context.Context, b *Bot, upd *Update, next func() error) error
```

Middleware wraps the dispatch of an update. Calling next runs the
remaining middlewares and, finally, the matching handler; not calling
it stops the update from reaching any handler.

### `Bot`

```go
type Bot struct
```

Bot is the bot API client. A Bot is safe for concurrent use; handlers,
middlewares and hooks may be registered or removed while it runs.

### `Stats`

```go
type Stats struct
```

Stats holds lightweight runtime counters.

### `StatsSnapshot`

```go
type StatsSnapshot struct
```

StatsSnapshot is a point-in-time copy of Stats.

### `New`

```go
func New(token string, opts ...Option) *Bot
```

New constructs a Bot using the supplied token. Options may be given
either as an Options value or through the With* helpers.

### `SetHTTPClient`

```go
func (b *Bot) SetHTTPClient(c *http.Client)
```

SetHTTPClient lets callers swap in a custom HTTP client.

### `Token`

```go
func (b *Bot) Token() string { return b.token }
```

Token returns the bot token.

### `BaseURL`

```go
func (b *Bot) BaseURL() string { return b.baseURL }
```

BaseURL returns the resolved bot API base URL.

### `Logger`

```go
func (b *Bot) Logger() Logger { return b.logger }
```

Logger returns the logger used by the bot.

### `Options`

```go
func (b *Bot) Options() Options { return b.opts }
```

Options returns a copy of the effective configuration.

### `Snapshot`

```go
func (b *Bot) Snapshot() StatsSnapshot
```

Snapshot returns the current runtime counters.

### `On`

```go
func (b *Bot) On(rules []Rule, handler Handler) int { return b.OnPriority(0, rules, handler) }
```

On registers a handler that fires when ALL of the supplied rules
match. The returned ID can be passed to Off.

### `PriorityPlugin`

```go
const PriorityPlugin = 50
```

PriorityPlugin is used for handlers registered by plugins, so a
catch-all handler of the application cannot swallow their commands.

### `PriorityFallback`

```go
const PriorityFallback = -100
```

PriorityFallback is used by Fallback: it always runs last.

### `OnPriority`

```go
func (b *Bot) OnPriority(prio int, rules []Rule, handler Handler) int
```

OnPriority is On with an explicit priority (default 0).

### `Off`

```go
func (b *Bot) Off(id int) bool
```

Off removes a handler registered with On. It reports whether the
handler existed.

### `OnStart`

```go
func (b *Bot) OnStart(h func(ctx context.Context, b *Bot) error) int
```

OnStart registers a startup hook and returns its ID.

### `OnStop`

```go
func (b *Bot) OnStop(h func(ctx context.Context, b *Bot) error) int
```

OnStop registers a shutdown hook and returns its ID.

### `RemoveHook`

```go
func (b *Bot) RemoveHook(id int) bool
```

RemoveHook removes a hook registered with OnStart or OnStop.

### `Use`

```go
func (b *Bot) Use(m Middleware) int
```

Use registers a middleware and returns its ID. Middlewares run in
registration order, outermost first.

### `RemoveMiddleware`

```go
func (b *Bot) RemoveMiddleware(id int) bool
```

RemoveMiddleware removes a middleware registered with Use.

### `Start`

```go
func (b *Bot) Start(ctx context.Context) error { return b.Run(ctx, "") }
```

Start runs the bot in polling mode. It blocks until ctx is
cancelled, Stop is called, or a fatal error occurs.

### `Run`

```go
func (b *Bot) Run(ctx context.Context, webhookURL string, hostAndPort ...string) error
```

Run starts the bot in polling mode, or in webhook mode when
webhookURL is non-empty (the optional argument is "host:port" of the
local listener).

### `IsRunning`

```go
func (b *Bot) IsRunning() bool { return b.running.Load() }
```

IsRunning reports whether the bot is currently running.

### `Stop`

```go
func (b *Bot) Stop()
```

Stop asks the polling loop or webhook server to terminate.

### `UserAgent`

```go
const UserAgent = core.UserAgent
```

Constants re-exported so callers don't have to import core.

### `Version2`

```go
const Version2 = core.Version
```

Constants re-exported so callers don't have to import core.
