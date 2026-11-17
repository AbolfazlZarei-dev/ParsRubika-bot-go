# `access` — Package access is a plugin that controls who may talk to the bot:

_Source: [`access.go`](../../../../ParsRubika/plugins/builtin/access/access.go)_

## Symbols

### `Plugin`

```go
type Plugin struct
```

Plugin implements plugins.Plugin and plugins.AccessService.

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

### `IsAdmin`

```go
func (p *Plugin) IsAdmin(id string) bool
```

IsAdmin implements plugins.AccessService.

### `IsBlocked`

```go
func (p *Plugin) IsBlocked(id string) bool
```

IsBlocked implements plugins.AccessService.

### `Block`

```go
func (p *Plugin) Block(id string) { p.mu.Lock(); p.blocked[id] = struct{}{}; p.mu.Unlock() }
```

Block adds a user to the block-list at runtime.

### `Unblock`

```go
func (p *Plugin) Unblock(id string) { p.mu.Lock(); delete(p.blocked, id); p.mu.Unlock() }
```

Unblock removes a user from the block-list.

### `AddAdmin`

```go
func (p *Plugin) AddAdmin(id string) { p.mu.Lock(); p.admins[id] = struct{}{}; p.mu.Unlock() }
```

AddAdmin grants administrator rights at runtime.
