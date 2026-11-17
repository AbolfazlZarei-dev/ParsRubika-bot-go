# `antilink` — Package antilink is a plugin for groups and channels: it deletes

_Source: [`antilink.go`](../../../../ParsRubika/plugins/builtin/antilink/antilink.go)_

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

### `Contains`

```go
func Contains(text string, allowedDomains []string, mentions bool) bool
```

Contains reports whether text holds a link (outside allowedDomains)
or, when mentions is set, an @mention.
