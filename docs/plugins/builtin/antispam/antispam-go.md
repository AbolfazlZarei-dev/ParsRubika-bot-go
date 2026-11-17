# `antispam` — Package antispam is a plugin that detects flooding and repeated

_Source: [`antispam.go`](../../../../ParsRubika/plugins/builtin/antispam/antispam.go)_

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

### `IsMuted`

```go
func (p *Plugin) IsMuted(uid string) bool
```

IsMuted reports whether a user is currently muted.

### `Unmute`

```go
func (p *Plugin) Unmute(uid string)
```

Unmute lifts a mute immediately.
