# `requestlog` — Package requestlog is a plugin that logs every update with its

_Source: [`requestlog.go`](../../../../ParsRubika/plugins/builtin/requestlog/requestlog.go)_

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
