# `plugintest` — Package plugintest provides a ready-made test harness for plugins: a

_Source: [`plugintest.go`](../../../ParsRubika/plugins/plugintest/plugintest.go)_

## Symbols

### `Call`

```go
type Call struct
```

Call is one recorded API request.

### `Responder`

```go
type Responder func(payload map[string]any) map[string]any
```

Responder customises the answer of the fake API for a method. It
returns the "data" object (nil means an empty object).

### `Env`

```go
type Env struct
```

Env is the test environment.

### `Printf`

```go
func (silent) Printf(string, ...any) {}
```

### `New`

```go
func New(t *testing.T, opts ...bot.Option) *Env
```

New builds an Env and registers its cleanup with t.

### `Respond`

```go
func (e *Env) Respond(method string, r Responder)
```

Respond installs a custom answer for an API method.

### `Calls`

```go
func (e *Env) Calls(method string) []Call
```

Calls returns the recorded calls of one method.

### `Reset`

```go
func (e *Env) Reset()
```

Reset forgets every recorded call.

### `Sent`

```go
func (e *Env) Sent() []string
```

Sent returns the texts of all sendMessage calls.

### `Enable`

```go
func (e *Env) Enable(names ...string)
```

Enable activates the named plugins and fails the test on error.

### `Say`

```go
func (e *Env) Say(chatID, senderID, text string) error
```

Say dispatches a text message from sender in chat.

### `Press`

```go
func (e *Env) Press(chatID, senderID, buttonID string) error
```

Press dispatches an inline button press.

### `Event`

```go
func (e *Env) Event(t models.UpdateType, chatID string) error
```

Event dispatches a non-message update (StartedBot, StoppedBot, ...).
