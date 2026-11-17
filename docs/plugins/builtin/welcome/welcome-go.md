# `welcome` — Package welcome is a plugin that greets users when they start the bot

_Source: [`welcome.go`](../../../../ParsRubika/plugins/builtin/welcome/welcome.go)_

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
