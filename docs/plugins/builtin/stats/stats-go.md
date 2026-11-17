# `stats` — Package stats is a plugin that exposes the bot's runtime counters

_Source: [`stats.go`](../../../../ParsRubika/plugins/builtin/stats/stats.go)_

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

### `Setup`

```go
func (p *Plugin) Setup(_ context.Context, h *plugins.Host) error
```

Setup implements plugins.Plugin.

### `Render`

```go
func Render(s bot.StatsSnapshot, plugins int) string
```

Render formats a snapshot for humans.
