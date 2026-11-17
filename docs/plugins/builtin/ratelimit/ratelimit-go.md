# `ratelimit` — Package ratelimit is a plugin that limits how fast each user may send

_Source: [`ratelimit.go`](../../../../ParsRubika/plugins/builtin/ratelimit/ratelimit.go)_

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

### `Validate`

```go
func (p *Plugin) Validate(c plugins.Config) error
```

Validate implements plugins.Validator.

### `Error`

```go
func (e errBad) Error() string { return string(e) }
```

### `Setup`

```go
func (p *Plugin) Setup(_ context.Context, h *plugins.Host) error
```

Setup implements plugins.Plugin.
