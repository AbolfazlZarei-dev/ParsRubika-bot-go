# `core` — Logger

_Source: [`logger.go`](../../ParsRubika/core/logger.go)_

## Symbols

### `Logger`

```go
type Logger interface
```

Logger is the minimal logging contract shared by every ParsRubika
component. *log.Logger and *logger.Logger both satisfy it.

### `LevelLogger`

```go
type LevelLogger interface
```

LevelLogger is the optional richer contract. When the installed
logger also implements it, components emit properly levelled
records instead of plain Printf lines.

### `Level`

```go
type Level int
```

Level identifies the severity used by Logf.

### `LogDebug`

```go
const LogDebug Level = iota
```

Severity levels understood by Logf.

### `LogInfo`

```go
const LogInfo
```

Severity levels understood by Logf.

### `LogWarn`

```go
const LogWarn
```

Severity levels understood by Logf.

### `LogError`

```go
const LogError
```

Severity levels understood by Logf.

### `Logf`

```go
func Logf(l Logger, lv Level, format string, args ...any)
```

Logf writes to l at the given level, degrading to Printf for loggers
that only implement the basic contract. A nil logger is ignored.
