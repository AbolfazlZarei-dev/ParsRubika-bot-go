# `bot` — Middleware

_Source: [`middleware.go`](../../ParsRubika/bot/middleware.go)_

## Symbols

### `TimerMiddleware`

```go
func TimerMiddleware(logFn func(format string, args ...any)) Middleware
```

TimerMiddleware logs the time taken to process each update.

### `RecovererMiddleware`

```go
func RecovererMiddleware(logFn func(format string, args ...any)) Middleware
```

RecovererMiddleware catches panics in downstream handlers/middleware
and logs the stack trace. Without this, a panic would crash the whole
process.

### `LoggingMiddleware`

```go
func LoggingMiddleware(logFn func(format string, args ...any)) Middleware
```

LoggingMiddleware logs every update at info level.

### `AllowedUsersMiddleware`

```go
func AllowedUsersMiddleware(allowedUserIDs []string, rejectMessage string) Middleware
```

AllowedUsersMiddleware only allows updates from the given user IDs.
Other updates are silently dropped (next() is not called).

### `BlockedUsersMiddleware`

```go
func BlockedUsersMiddleware(blockedUserIDs []string) Middleware
```

BlockedUsersMiddleware drops updates from the given user IDs.

### `PerUserRateLimitMiddleware`

```go
func PerUserRateLimitMiddleware(maxUpdates int, window time.Duration) Middleware
```

PerUserRateLimitMiddleware limits how many updates per user per
window are processed. Excess updates are silently dropped.

### `ChatTypeMiddleware`

```go
func ChatTypeMiddleware(chatTypes ...string) Middleware
```

ChatTypeMiddleware only allows updates from the given chat types
("private", "group", "channel").

### `CommandMiddleware`

```go
func CommandMiddleware(prefixes ...string) Middleware
```

CommandMiddleware extracts the command from the message text (if it
starts with a prefix like / or !) and stores it in ExternalData for
downstream rules to use.

### `UserStatsMiddleware`

```go
func UserStatsMiddleware(tracker *MetricsTracker) Middleware
```

UserStatsMiddleware tracks per-user message counts in a MetricsTracker.

### `ChainMiddlewares`

```go
func ChainMiddlewares(mws ...Middleware) Middleware
```

ChainMiddlewares combines multiple middlewares into a single middleware
that runs them in order. This is useful when you want to apply a group
of middlewares as a unit.

### `EnsureUpdate`

```go
func EnsureUpdate(upd *Update) *models.Update
```

EnsureUpdate is a helper that returns the *models.Update from a bot
update, or nil if the update is nil. This is useful in middleware
that needs to access update fields.
