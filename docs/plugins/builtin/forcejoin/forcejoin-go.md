# `forcejoin` — Package forcejoin is a plugin that makes users join one or more

_Source: [`forcejoin.go`](../../../../ParsRubika/plugins/builtin/forcejoin/forcejoin.go)_

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

### `Error`

```go
func (e constErr) Error() string { return string(e) }
```

### `Validate`

```go
func (p *Plugin) Validate(c plugins.Config) error
```

Validate implements plugins.Validator.

### `Setup`

```go
func (p *Plugin) Setup(_ context.Context, h *plugins.Host) error
```

Setup implements plugins.Plugin.
