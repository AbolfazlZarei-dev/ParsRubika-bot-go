# `users` — Package users is a plugin that remembers every chat the bot talks to

_Source: [`users.go`](../../../../ParsRubika/plugins/builtin/users/users.go)_

## Symbols

### `Record`

```go
type Record struct
```

Record is what is known about one chat.

### `Plugin`

```go
type Plugin struct
```

Plugin implements plugins.Plugin and plugins.UserDirectory.

### `New`

```go
func New() *Plugin
```

New creates the plugin.

### `Setup`

```go
func (p *Plugin) Setup(_ context.Context, h *plugins.Host) error
```

Setup implements plugins.Plugin.

### `Teardown`

```go
func (p *Plugin) Teardown(context.Context) error { return p.Flush() }
```

Teardown implements plugins.Teardowner.

### `Flush`

```go
func (p *Plugin) Flush() error
```

Flush writes pending changes to disk.

### `Chats`

```go
func (p *Plugin) Chats() []string
```

Chats implements plugins.UserDirectory: active chats, sorted.

### `Count`

```go
func (p *Plugin) Count() int { return len(p.Chats()) }
```

Count implements plugins.UserDirectory.

### `Get`

```go
func (p *Plugin) Get(chatID string) (Record, bool)
```

Get returns the record of a chat.
