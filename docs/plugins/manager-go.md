# `plugins` — Manager

_Source: [`manager.go`](../../ParsRubika/plugins/manager.go)_

## Symbols

### `Manager`

```go
type Manager struct
```

Manager registers plugins and controls their life cycle.

### `Option`

```go
type Option func(*Manager)
```

Option customises a Manager.

### `WithConfigs`

```go
func WithConfigs(c map[string]Config) Option
```

WithConfigs supplies per-plugin configuration (typically loaded with
LoadConfigFile). Entries apply to plugins registered later too.

### `New`

```go
func New(b *bot.Bot, opts ...Option) *Manager
```

New creates a Manager bound to b.

### `Bot`

```go
func (m *Manager) Bot() *bot.Bot { return m.bot }
```

Bot returns the bot the manager is bound to.

### `Register`

```go
func (m *Manager) Register(p Plugin, cfg ...Config) error
```

Register adds a plugin. cfg (optional) overrides the plugin defaults
and any configuration supplied through WithConfigs.

### `MustRegister`

```go
func (m *Manager) MustRegister(p Plugin, cfg ...Config) *Manager
```

MustRegister is Register that panics on error; convenient in main().

### `Unregister`

```go
func (m *Manager) Unregister(ctx context.Context, name string) error
```

Unregister disables and removes a plugin.

### `Configure`

```go
func (m *Manager) Configure(name string, cfg Config) error
```

Configure replaces the user configuration of a plugin. It takes
effect the next time the plugin is enabled (or with Reload).

### `Enable`

```go
func (m *Manager) Enable(ctx context.Context, names ...string) error
```

Enable activates the plugins (and, first, everything they require).

### `EnableAll`

```go
func (m *Manager) EnableAll(ctx context.Context) error
```

EnableAll activates every registered plugin that is not marked
Disabled, ordered by priority and then registration order.

### `Disable`

```go
func (m *Manager) Disable(ctx context.Context, names ...string) error
```

Disable deactivates a plugin. Plugins that require it are disabled
first. Everything the plugin registered is removed from the bot.

### `DisableAll`

```go
func (m *Manager) DisableAll(ctx context.Context) error
```

DisableAll deactivates every plugin, newest first.

### `Reload`

```go
func (m *Manager) Reload(ctx context.Context, name string, cfg ...Config) error
```

Reload disables and re-enables a plugin, optionally with a new
configuration.

### `Install`

```go
func (m *Manager) Install()
```

Install hooks the manager into the bot life cycle: EnableAll runs when
the bot starts and DisableAll when it stops.

### `Info`

```go
type Info struct
```

Info is a read-only snapshot of a plugin.

### `List`

```go
func (m *Manager) List() []Info
```

List returns every registered plugin, ordered by name.

### `Get`

```go
func (m *Manager) Get(name string) (Plugin, bool)
```

Get returns a registered plugin by name.

### `IsEnabled`

```go
func (m *Manager) IsEnabled(name string) bool
```

IsEnabled reports whether a plugin is currently active.

### `Enabled`

```go
func (m *Manager) Enabled() []string
```

Enabled returns the names of active plugins in activation order.

### `HelpEntries`

```go
func (m *Manager) HelpEntries() []HelpEntry
```

HelpEntries returns the documented commands of all active plugins.

### `Service`

```go
func (m *Manager) Service(key string) (any, bool)
```

Service returns a published service.
