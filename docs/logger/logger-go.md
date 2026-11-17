# `logger` — Package logger is a small, dependency-free, levelled logger with

_Source: [`logger.go`](../../ParsRubika/logger/logger.go)_

## Symbols

### `Field`

```go
type Field struct
```

Field is a structured key/value pair attached to a record.

### `F`

```go
func F(key string, value any) Field { return Field{Key: key, Value: value} }
```

F builds a Field.

### `Err`

```go
func Err(err error) Field { return Field{Key: "error", Value: err} }
```

Err builds an "error" field. A nil error yields an empty field that
is skipped by the formatters.

### `Entry`

```go
type Entry struct
```

Entry is a fully assembled log record, handed to hooks.

### `Format`

```go
type Format int
```

Format selects the textual layout of a sink.

### `FormatText`

```go
const FormatText Format = iota
```

FormatText is the human friendly single-line layout.

### `FormatJSON`

```go
const FormatJSON
```

FormatJSON emits one JSON object per line.

### `Logger`

```go
type Logger struct
```

Logger writes records to its sinks. Child loggers created with With
or Named share the sinks and level of their parent. A Logger is safe
for concurrent use.

### `Option`

```go
type Option func(*config)
```

Option configures a Logger created by New.

### `WithLevel`

```go
func WithLevel(l Level) Option { return func(c *config) { c.level = l } }
```

WithLevel sets the minimum level of the logger.

### `WithName`

```go
func WithName(n string) Option { return func(c *config) { c.name = n } }
```

WithName sets the logger name shown in every record.

### `WithFields`

```go
func WithFields(f ...Field) Option { return func(c *config) { c.fields = append(c.fields, f...) } }
```

WithFields attaches permanent fields to every record.

### `WithTimeFormat`

```go
func WithTimeFormat(layout string) Option { return func(c *config) { c.timeFmt = layout } }
```

WithTimeFormat changes the timestamp layout ("" hides the timestamp).

### `WithCaller`

```go
func WithCaller(on bool) Option { return func(c *config) { c.caller = on } }
```

WithCaller adds the call site (file:line) to every record.

### `WithHook`

```go
func WithHook(h func(Entry)) Option { return func(c *config) { c.hooks = append(c.hooks, h) } }
```

WithHook registers a callback invoked for every emitted record.

### `WithOutput`

```go
func WithOutput(w io.Writer, mode ColorMode) Option
```

WithOutput replaces the default stderr sink with w using text format.
Colours follow mode (ColorAuto inspects w).

### `WithColor`

```go
func WithColor(mode ColorMode) Option
```

WithColor overrides the colour mode of the default console sink.

### `WithJSON`

```go
func WithJSON() Option
```

WithJSON switches the default console sink to JSON lines.

### `AddOutput`

```go
func AddOutput(w io.Writer, format Format, mode ColorMode, min Level) Option
```

AddOutput attaches an additional sink. Records below min are skipped
by this sink only.

### `WithFile`

```go
func WithFile(path string, format Format, maxSizeMB, maxBackups int) Option
```

WithFile adds a plain-text (or JSON) rotating log file sink.
maxSizeMB <= 0 disables rotation.

### `New`

```go
func New(opts ...Option) *Logger
```

New builds a Logger. Without options it writes coloured text to
stderr at Info level (or at the level named by PARSRUBIKA_LOG_LEVEL).

### `Nop`

```go
func Nop() *Logger { return New(WithLevel(LevelOff), WithOutput(io.Discard, ColorNever)) }
```

Nop returns a logger that discards everything.

### `Default`

```go
func Default() *Logger
```

Default returns the process-wide logger, creating it on first use.

### `SetDefault`

```go
func SetDefault(l *Logger)
```

SetDefault replaces the process-wide logger.

### `With`

```go
func (l *Logger) With(fields ...Field) *Logger
```

With returns a child logger carrying extra permanent fields.

### `Named`

```go
func (l *Logger) Named(name string) *Logger
```

Named returns a child logger with a (dot-joined) component name.

### `SetLevel`

```go
func (l *Logger) SetLevel(lv Level) { l.c.level.Store(int32(lv)) }
```

SetLevel changes the minimum level at runtime (affects all children).

### `Level`

```go
func (l *Logger) Level() Level { return Level(l.c.level.Load()) }
```

Level returns the current minimum level.

### `Enabled`

```go
func (l *Logger) Enabled(lv Level) bool { return lv >= l.Level() && lv < LevelOff }
```

Enabled reports whether a record of the given level would be emitted.

### `Close`

```go
func (l *Logger) Close() error
```

Close flushes and closes every file sink.

### `Trace`

```go
func (l *Logger) Trace(msg string, f ...Field) { l.log(LevelTrace, msg, f) }
```

### `Debug`

```go
func (l *Logger) Debug(msg string, f ...Field) { l.log(LevelDebug, msg, f) }
```

### `Info`

```go
func (l *Logger) Info(msg string, f ...Field) { l.log(LevelInfo, msg, f) }
```

### `Warn`

```go
func (l *Logger) Warn(msg string, f ...Field) { l.log(LevelWarn, msg, f) }
```

### `Error`

```go
func (l *Logger) Error(msg string, f ...Field) { l.log(LevelError, msg, f) }
```

### `Tracef`

```go
func (l *Logger) Tracef(format string, a ...any) { l.logf(LevelTrace, format, a) }
```

### `Debugf`

```go
func (l *Logger) Debugf(format string, a ...any) { l.logf(LevelDebug, format, a) }
```

### `Infof`

```go
func (l *Logger) Infof(format string, a ...any) { l.logf(LevelInfo, format, a) }
```

### `Warnf`

```go
func (l *Logger) Warnf(format string, a ...any) { l.logf(LevelWarn, format, a) }
```

### `Errorf`

```go
func (l *Logger) Errorf(format string, a ...any) { l.logf(LevelError, format, a) }
```

### `Printf`

```go
func (l *Logger) Printf(format string, a ...any) { l.logf(LevelInfo, format, a) }
```

Printf makes *Logger usable wherever a standard-library style logger
is expected. The record is emitted at Info level.

### `Writer`

```go
func (l *Logger) Writer(lv Level) io.Writer { return &lineWriter{l: l, lv: lv} }
```

Writer returns an io.Writer that logs every written line at lv.

### `StdLogger`

```go
func (l *Logger) StdLogger(lv Level) *log.Logger { return log.New(l.Writer(lv), "", 0) }
```

StdLogger adapts the logger to the standard library type, for APIs
that insist on *log.Logger.

### `Write`

```go
func (w *lineWriter) Write(p []byte) (int, error)
```

### `Trace`

```go
func Trace(msg string, f ...Field) { Default().log(LevelTrace, msg, f) }
```

### `Debug`

```go
func Debug(msg string, f ...Field) { Default().log(LevelDebug, msg, f) }
```

### `Info`

```go
func Info(msg string, f ...Field) { Default().log(LevelInfo, msg, f) }
```

### `Warn`

```go
func Warn(msg string, f ...Field) { Default().log(LevelWarn, msg, f) }
```

### `Error`

```go
func Error(msg string, f ...Field) { Default().log(LevelError, msg, f) }
```

### `Infof`

```go
func Infof(format string, a ...any) { Default().logf(LevelInfo, format, a) }
```

### `Warnf`

```go
func Warnf(format string, a ...any) { Default().logf(LevelWarn, format, a) }
```

### `Errorf`

```go
func Errorf(format string, a ...any) { Default().logf(LevelError, format, a) }
```

### `Debugf`

```go
func Debugf(format string, a ...any) { Default().logf(LevelDebug, format, a) }
```
