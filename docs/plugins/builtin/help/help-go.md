# `help` — Package help is a plugin that answers /help with the commands that

_Source: [`help.go`](../../../../ParsRubika/plugins/builtin/help/help.go)_

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

### `DefaultConfig`

```go
func (p *Plugin) DefaultConfig() plugins.Config
```

DefaultConfig implements plugins.Defaulter.

### `Setup`

```go
func (p *Plugin) Setup(_ context.Context, h *plugins.Host) error
```

Setup implements plugins.Plugin.

### `Render`

```go
func Render(m *plugins.Manager, header, footer string) string
```

Render builds the help text from the documented commands of all
active plugins.
