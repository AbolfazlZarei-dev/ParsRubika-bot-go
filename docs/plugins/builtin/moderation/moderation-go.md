# `moderation` — Package moderation is a plugin with group moderation commands for

_Source: [`moderation.go`](../../../../ParsRubika/plugins/builtin/moderation/moderation.go)_

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
func (p *Plugin) DefaultConfig() plugins.Config { return plugins.Config{"max_warns": 3} }
```

DefaultConfig implements plugins.Defaulter.

### `Setup`

```go
func (p *Plugin) Setup(_ context.Context, h *plugins.Host) error
```

Setup implements plugins.Plugin.
