# `builtin` — Package builtin gathers the ready-made plugins that ship with

_Source: [`builtin.go`](../../../ParsRubika/plugins/builtin/builtin.go)_

## Symbols

### `All`

```go
func All() []plugins.Plugin
```

All returns a fresh instance of every built-in plugin.

### `RegisterAll`

```go
func RegisterAll(m *plugins.Manager, cfg map[string]plugins.Config) error
```

RegisterAll registers every built-in plugin on m. cfg holds optional
per-plugin configuration keyed by plugin name. Plugins that cannot
work without configuration (forcejoin) are only enabled when their
configuration is present.
