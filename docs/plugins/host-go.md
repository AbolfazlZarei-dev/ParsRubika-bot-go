# `plugins` — Host

_Source: [`host.go`](../../ParsRubika/plugins/host.go)_

## Symbols

### `HelpEntry`

```go
type HelpEntry struct
```

HelpEntry documents one user-visible command of a plugin.

### `Host`

```go
type Host struct
```

Host is the plugin's window into the bot. Everything registered
through it is tracked and removed again when the plugin is disabled,
so plugins never have to clean up after themselves.

### `Name`

```go
func (h *Host) Name() string { return h.name }
```

Name returns the plugin name.

### `Bot`

```go
func (h *Host) Bot() *bot.Bot { return h.bot }
```

Bot returns the bot the plugin is attached to.

### `Manager`

```go
func (h *Host) Manager() *Manager { return h.m }
```

Manager returns the manager that owns the plugin.

### `Config`

```go
func (h *Host) Config() Config { return h.cfg }
```

Config returns the effective configuration (defaults merged with the
values supplied at registration time).

### `Context`

```go
func (h *Host) Context() context.Context { return h.ctx }
```

Context is cancelled when the plugin is disabled. Long-running
goroutines of the plugin should stop when it is done.

### `Debugf`

```go
func (h *Host) Debugf(format string, a ...any) { h.logf(core.LogDebug, format, a...) }
```

Debugf logs at debug level, prefixed with the plugin name.

### `Infof`

```go
func (h *Host) Infof(format string, a ...any) { h.logf(core.LogInfo, format, a...) }
```

Infof logs at info level, prefixed with the plugin name.

### `Warnf`

```go
func (h *Host) Warnf(format string, a ...any) { h.logf(core.LogWarn, format, a...) }
```

Warnf logs at warning level, prefixed with the plugin name.

### `Errorf`

```go
func (h *Host) Errorf(format string, a ...any) { h.logf(core.LogError, format, a...) }
```

Errorf logs at error level, prefixed with the plugin name.

### `On`

```go
func (h *Host) On(rs []bot.Rule, fn bot.Handler) int
```

On registers a low-level handler. Plugin handlers run before the
application's ordinary handlers (see bot.PriorityPlugin).

### `Handle`

```go
func (h *Host) Handle(rs []bot.Rule, fn bot.HandlerFunc) int
```

Handle registers an ergonomic handler for updates matching all rules.

### `Command`

```go
func (h *Host) Command(name string, fn bot.HandlerFunc, aliases ...string) int
```

Command registers a "/name" (or "!name") handler plus aliases.

### `Text`

```go
func (h *Host) Text(fn bot.HandlerFunc, exact ...string) int
```

Text registers a handler for exact message texts.

### `Regex`

```go
func (h *Host) Regex(pattern string, fn bot.HandlerFunc) int
```

Regex registers a handler for texts matching pattern.

### `Button`

```go
func (h *Host) Button(id string, fn bot.HandlerFunc) int
```

Button registers a handler for an inline button ID.

### `Message`

```go
func (h *Host) Message(fn bot.HandlerFunc) int
```

Message registers a handler for every text message.

### `Use`

```go
func (h *Host) Use(mw bot.Middleware) int { return h.track(&h.mws, h.bot.Use(mw)) }
```

Use registers a middleware.

### `OnStart`

```go
func (h *Host) OnStart(fn func(ctx context.Context, b *bot.Bot) error) int
```

OnStart registers a bot start-up hook.

### `OnStop`

```go
func (h *Host) OnStop(fn func(ctx context.Context, b *bot.Bot) error) int
```

OnStop registers a bot shutdown hook.

### `Help`

```go
func (h *Host) Help(command, description string)
```

Help documents a command so the help plugin can list it.

### `Cleanup`

```go
func (h *Host) Cleanup(fn func())
```

Cleanup schedules fn to run (last in, first out) when the plugin is
disabled.

### `Every`

```go
func (h *Host) Every(d time.Duration, fn func(ctx context.Context))
```

Every runs fn repeatedly with the given period until the plugin is
disabled. A panic in fn is recovered and logged.

### `After`

```go
func (h *Host) After(d time.Duration, fn func(ctx context.Context))
```

After runs fn once after d unless the plugin is disabled first.

### `Go`

```go
func (h *Host) Go(fn func(ctx context.Context))
```

Go starts fn in a goroutine that is awaited when the plugin is disabled.

### `Provide`

```go
func (h *Host) Provide(key string, svc any) error
```

Provide publishes a service other plugins (and the application) can
look up with Service. The key must be unique; the service is
withdrawn when the plugin is disabled.

### `Service`

```go
func (h *Host) Service(key string) (any, bool) { return h.m.Service(key) }
```

Service looks up a service published by another plugin.
