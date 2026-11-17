# `plugins` — Config

_Source: [`config.go`](../../ParsRubika/plugins/config.go)_

## Symbols

### `Config`

```go
type Config map[string]any
```

Config is the free-form configuration of one plugin. Values coming
from JSON files are float64/[]any/map[string]any; the typed getters
understand both those and native Go values.

### `Clone`

```go
func (c Config) Clone() Config
```

Clone returns a deep-enough copy (nested maps and slices are copied).

### `Merge`

```go
func Merge(defaults, override Config) Config
```

Merge returns defaults overlaid with override.

### `String`

```go
func (c Config) String(key, def string) string
```

String returns the value of key as a string.

### `Int`

```go
func (c Config) Int(key string, def int) int
```

Int returns the value of key as an int.

### `Int64`

```go
func (c Config) Int64(key string, def int64) int64
```

Int64 returns the value of key as an int64.

### `Bool`

```go
func (c Config) Bool(key string, def bool) bool
```

Bool returns the value of key as a bool.

### `Duration`

```go
func (c Config) Duration(key string, def time.Duration) time.Duration
```

Duration returns the value of key as a time.Duration. Strings such as
"10s" or "2m" are parsed; plain numbers are read as seconds.

### `Strings`

```go
func (c Config) Strings(key string) []string
```

Strings returns the value of key as a string slice.

### `Has`

```go
func (c Config) Has(key string) bool { _, ok := c[key]; return ok }
```

Has reports whether key is present.

### `ParseConfigJSON`

```go
func ParseConfigJSON(data []byte) (map[string]Config, error)
```

ParseConfigJSON parses a document of the form

{"plugins": {"antispam": {"max_messages": 5}, "welcome": {...}}}

into per-plugin configurations (keys are lower-cased).

### `LoadConfigFile`

```go
func LoadConfigFile(path string) (map[string]Config, error)
```

LoadConfigFile reads a JSON configuration file (see ParseConfigJSON).
