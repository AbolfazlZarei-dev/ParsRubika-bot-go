# `kv` — Package kv is a plugin that publishes a small key/value store with

_Source: [`kv.go`](../../../../ParsRubika/plugins/builtin/kv/kv.go)_

## Symbols

### `Store`

```go
type Store interface
```

Store is the published service.

### `Plugin`

```go
type Plugin struct
```

Plugin implements plugins.Plugin and Store.

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

### `Teardown`

```go
func (p *Plugin) Teardown(context.Context) error { return p.Flush() }
```

Teardown implements plugins.Teardowner: a final flush.

### `Flush`

```go
func (p *Plugin) Flush() error
```

Flush writes pending changes to disk.

### `Get`

```go
func (p *Plugin) Get(key string) (any, bool)
```

Get implements Store.

### `Set`

```go
func (p *Plugin) Set(key string, value any) { p.SetTTL(key, value, 0) }
```

Set implements Store.

### `SetTTL`

```go
func (p *Plugin) SetTTL(key string, value any, ttl time.Duration)
```

SetTTL implements Store; ttl <= 0 means no expiry.

### `Delete`

```go
func (p *Plugin) Delete(key string)
```

Delete implements Store.

### `Keys`

```go
func (p *Plugin) Keys(prefix string) []string
```

Keys implements Store: live keys with the prefix, sorted.
