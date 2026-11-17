# `plugins` — Package plugins is the plugin framework of ParsRubika

_Source: [`plugin.go`](../../ParsRubika/plugins/plugin.go)_

## Symbols

### `DefaultPriority`

```go
const DefaultPriority = 100
```

DefaultPriority is the start-up priority of plugins that do not set one.
Plugins with a lower number are enabled earlier.

### `Meta`

```go
type Meta struct
```

Meta describes a plugin.

### `Plugin`

```go
type Plugin interface
```

Plugin is the contract every plugin fulfils.

### `Teardowner`

```go
type Teardowner interface
```

Teardowner is implemented by plugins that need to release resources of
their own (files, connections) when they are disabled.

### `Defaulter`

```go
type Defaulter interface
```

Defaulter lets a plugin declare its default configuration. User
supplied values are merged on top of it.

### `Validator`

```go
type Validator interface
```

Validator lets a plugin reject an invalid configuration before Setup.

### `Base`

```go
type Base struct
```

Base is a convenience type to embed in plugins: it implements Meta.

### `NewBase`

```go
func NewBase(m Meta) Base { return Base{meta: m} }
```

NewBase builds a Base from m.

### `Meta`

```go
func (b Base) Meta() Meta { return b.meta }
```

Meta implements Plugin.

### `ErrInvalidName`

```go
var ErrInvalidName = errors.New("plugins: invalid plugin name")
```

Errors returned by the Manager.

### `ErrDuplicate`

```go
var ErrDuplicate = errors.New("plugins: plugin already registered")
```

Errors returned by the Manager.

### `ErrNotFound`

```go
var ErrNotFound = errors.New("plugins: plugin not registered")
```

Errors returned by the Manager.

### `ErrDependency`

```go
var ErrDependency = errors.New("plugins: dependency problem")
```

Errors returned by the Manager.

### `ErrNoBot`

```go
var ErrNoBot = errors.New("plugins: manager has no bot")
```

Errors returned by the Manager.

### `ErrAlreadyActive`

```go
var ErrAlreadyActive = errors.New("plugins: plugin already enabled")
```

Errors returned by the Manager.

### `SetupError`

```go
type SetupError struct
```

SetupError wraps a failure raised while enabling a plugin.

### `Error`

```go
func (e *SetupError) Error() string { return fmt.Sprintf("plugins: enable %q: %v", e.Plugin, e.Err) }
```

Error implements error.

### `Unwrap`

```go
func (e *SetupError) Unwrap() error { return e.Err }
```

Unwrap returns the underlying error.
