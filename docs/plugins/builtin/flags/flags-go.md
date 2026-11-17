# `flags` — Package flags is a plugin that publishes feature flags (the "flags"

_Source: [`flags.go`](../../../../ParsRubika/plugins/builtin/flags/flags.go)_

## Symbols

### `Flags`

```go
type Flags interface
```

Flags is the published service.

### `Plugin`

```go
type Plugin struct
```

Plugin implements plugins.Plugin and Flags.

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

### `IsEnabled`

```go
func (p *Plugin) IsEnabled(name string) bool
```

IsEnabled implements Flags.

### `IsEnabledFor`

```go
func (p *Plugin) IsEnabledFor(name, userID string) bool
```

IsEnabledFor implements Flags: a per-user override wins over the
global value.

### `Set`

```go
func (p *Plugin) Set(name string, on bool) { p.mu.Lock(); p.global[name] = on; p.mu.Unlock() }
```

Set implements Flags.

### `SetFor`

```go
func (p *Plugin) SetFor(name, userID string, on bool)
```

SetFor implements Flags.
