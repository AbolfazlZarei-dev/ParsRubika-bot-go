# `logger` — Level

_Source: [`level.go`](../../ParsRubika/logger/level.go)_

## Symbols

### `Level`

```go
type Level int32
```

Level is the severity of a log record.

### `LevelTrace`

```go
const LevelTrace Level = iota - 1
```

Severity levels, from the most verbose to the most critical.

### `LevelDebug`

```go
const LevelDebug
```

Severity levels, from the most verbose to the most critical.

### `LevelInfo`

```go
const LevelInfo
```

Severity levels, from the most verbose to the most critical.

### `LevelWarn`

```go
const LevelWarn
```

Severity levels, from the most verbose to the most critical.

### `LevelError`

```go
const LevelError
```

Severity levels, from the most verbose to the most critical.

### `LevelOff`

```go
const LevelOff
```

Severity levels, from the most verbose to the most critical.

### `String`

```go
func (l Level) String() string
```

String returns the upper-case name of the level.

### `ParseLevel`

```go
func ParseLevel(s string) (Level, error)
```

ParseLevel converts a textual level ("debug", "warn", ...) to a Level.
