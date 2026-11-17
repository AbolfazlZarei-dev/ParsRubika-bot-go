# `router` — Package router provides a simple request router for the bot API

_Source: [`router.go`](../../../ParsRubika/bot/router/router.go)_

## Symbols

### `Handler`

```go
type Handler func(ctx context.Context, b *bot.Bot, upd *models.Update) error
```

Handler is the signature of a routed handler.

### `Route`

```go
type Route struct
```

Route is a single routing entry: a set of rules and a handler.

### `Router`

```go
type Router struct
```

Router holds a list of routes and dispatches updates to the first
matching one.

### `New`

```go
func New() *Router { return &Router{} }
```

New constructs an empty Router.

### `Add`

```go
func (r *Router) Add(rs []rules.Rule, h Handler) *Router
```

Add appends a route to the router.

### `Command`

```go
func (r *Router) Command(cmds []string, h Handler) *Router
```

Command adds a route that matches a specific command.

### `Text`

```go
func (r *Router) Text(text string, h Handler) *Router
```

Text adds a route that matches a text message.

### `Regex`

```go
func (r *Router) Regex(pattern string, h Handler) *Router
```

Regex adds a route that matches a regex against the text.

### `AnyText`

```go
func (r *Router) AnyText(h Handler) *Router
```

AnyText adds a route that matches any text message.

### `Button`

```go
func (r *Router) Button(buttonID string, h Handler) *Router
```

Button adds a route that matches a button click by button_id.

### `Group`

```go
func (r *Router) Group(h Handler) *Router
```

Group adds a route that matches updates from groups.

### `Private`

```go
func (r *Router) Private(h Handler) *Router
```

Private adds a route that matches updates from private chats.

### `Channel`

```go
func (r *Router) Channel(h Handler) *Router
```

Channel adds a route that matches updates from channels.

### `Attach`

```go
func (r *Router) Attach(b *bot.Bot)
```

Attach registers the router on a bot. The first matching route's
handler is invoked; subsequent routes are skipped (first-match
semantics).

### `PrintRoutes`

```go
func (r *Router) PrintRoutes() string
```

PrintRoutes returns a human-readable list of all registered
routes (for debugging).
