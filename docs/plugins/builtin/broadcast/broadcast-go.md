# `broadcast` — Package broadcast is a plugin that lets administrators send one

_Source: [`broadcast.go`](../../../../ParsRubika/plugins/builtin/broadcast/broadcast.go)_

## Symbols

### `Plugin`

```go
type Plugin struct
```

Plugin implements plugins.Plugin.

### `New`

```go
func New() *Plugin
```

New creates the plugin.

### `Result`

```go
type Result struct
```

Result summarises a finished broadcast.

### `Setup`

```go
func (p *Plugin) Setup(_ context.Context, h *plugins.Host) error
```

Setup implements plugins.Plugin.

### `Send`

```go
func Send(ctx context.Context, b *bot.Bot, chats []string, text string, delay time.Duration) Result
```

Send delivers text to every chat sequentially, pausing delay between
calls, and stops early when ctx ends.
