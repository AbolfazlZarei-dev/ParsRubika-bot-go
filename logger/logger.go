// Package logger is a small, dependency-free, levelled logger with
// coloured console output, structured fields, JSON mode and rotating
// file sinks.
//
// It is entirely optional: every ParsRubika component accepts any value
// with a Printf method, so the standard library logger works as well.
// Plug this one in when you want readable, colourful and structured
// logs:
//
//	log := logger.New(logger.WithLevel(logger.LevelDebug))
//	b := bot.New(token, bot.WithLogger(log))
//	log.Info("ready", logger.F("user", 42))
package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// Field is a structured key/value pair attached to a record.
type Field struct {
	Key   string
	Value any
}

// F builds a Field.
func F(key string, value any) Field { return Field{Key: key, Value: value} }

// Err builds an "error" field. A nil error yields an empty field that
// is skipped by the formatters.
func Err(err error) Field { return Field{Key: "error", Value: err} }

// Entry is a fully assembled log record, handed to hooks.
type Entry struct {
	Time    time.Time
	Level   Level
	Name    string
	Message string
	Fields  []Field
	Caller  string
}

// Format selects the textual layout of a sink.
type Format int

const (
	// FormatText is the human friendly single-line layout.
	FormatText Format = iota
	// FormatJSON emits one JSON object per line.
	FormatJSON
)

type sink struct {
	w      io.Writer
	closer io.Closer
	format Format
	color  bool
	min    Level
}

type core struct {
	level   atomic.Int32
	mu      sync.Mutex
	sinks   []*sink
	timeFmt string
	caller  bool
	hooks   []func(Entry)
	pool    sync.Pool
}

// Logger writes records to its sinks. Child loggers created with With
// or Named share the sinks and level of their parent. A Logger is safe
// for concurrent use.
type Logger struct {
	c      *core
	name   string
	fields []Field
}

// Option configures a Logger created by New.
type Option func(*config)

type config struct {
	level   Level
	timeFmt string
	caller  bool
	name    string
	fields  []Field
	sinks   []sinkSpec
	hooks   []func(Entry)
}

type sinkSpec struct {
	w      io.Writer
	closer io.Closer
	format Format
	mode   ColorMode
	min    Level
}

// WithLevel sets the minimum level of the logger.
func WithLevel(l Level) Option { return func(c *config) { c.level = l } }

// WithName sets the logger name shown in every record.
func WithName(n string) Option { return func(c *config) { c.name = n } }

// WithFields attaches permanent fields to every record.
func WithFields(f ...Field) Option { return func(c *config) { c.fields = append(c.fields, f...) } }

// WithTimeFormat changes the timestamp layout ("" hides the timestamp).
func WithTimeFormat(layout string) Option { return func(c *config) { c.timeFmt = layout } }

// WithCaller adds the call site (file:line) to every record.
func WithCaller(on bool) Option { return func(c *config) { c.caller = on } }

// WithHook registers a callback invoked for every emitted record.
func WithHook(h func(Entry)) Option { return func(c *config) { c.hooks = append(c.hooks, h) } }

// WithOutput replaces the default stderr sink with w using text format.
// Colours follow mode (ColorAuto inspects w).
func WithOutput(w io.Writer, mode ColorMode) Option {
	return func(c *config) { c.sinks = []sinkSpec{{w: w, format: FormatText, mode: mode, min: LevelTrace}} }
}

// WithColor overrides the colour mode of the default console sink.
func WithColor(mode ColorMode) Option {
	return func(c *config) {
		if len(c.sinks) == 1 && c.sinks[0].w == os.Stderr {
			c.sinks[0].mode = mode
		}
	}
}

// WithJSON switches the default console sink to JSON lines.
func WithJSON() Option {
	return func(c *config) {
		if len(c.sinks) == 1 && c.sinks[0].w == os.Stderr {
			c.sinks[0].format = FormatJSON
		}
	}
}

// AddOutput attaches an additional sink. Records below min are skipped
// by this sink only.
func AddOutput(w io.Writer, format Format, mode ColorMode, min Level) Option {
	return func(c *config) {
		c.sinks = append(c.sinks, sinkSpec{w: w, format: format, mode: mode, min: min})
	}
}

// WithFile adds a plain-text (or JSON) rotating log file sink.
// maxSizeMB <= 0 disables rotation.
func WithFile(path string, format Format, maxSizeMB, maxBackups int) Option {
	return func(c *config) {
		rf, err := NewRotatingFile(path, maxSizeMB, maxBackups)
		if err != nil {
			fmt.Fprintf(os.Stderr, "logger: cannot open %s: %v\n", path, err)
			return
		}
		c.sinks = append(c.sinks, sinkSpec{w: rf, closer: rf, format: format, mode: ColorNever, min: LevelTrace})
	}
}

// New builds a Logger. Without options it writes coloured text to
// stderr at Info level (or at the level named by PARSRUBIKA_LOG_LEVEL).
func New(opts ...Option) *Logger {
	cfg := config{
		level:   LevelInfo,
		timeFmt: "15:04:05.000",
		sinks:   []sinkSpec{{w: os.Stderr, format: FormatText, mode: ColorAuto, min: LevelTrace}},
	}
	if v := os.Getenv("PARSRUBIKA_LOG_LEVEL"); v != "" {
		if lv, err := ParseLevel(v); err == nil {
			cfg.level = lv
		}
	}
	for _, o := range opts {
		o(&cfg)
	}
	c := &core{timeFmt: cfg.timeFmt, caller: cfg.caller, hooks: cfg.hooks}
	c.level.Store(int32(cfg.level))
	c.pool.New = func() any { b := make([]byte, 0, 256); return &b }
	for _, s := range cfg.sinks {
		c.sinks = append(c.sinks, &sink{
			w: s.w, closer: s.closer, format: s.format, min: s.min,
			color: s.format == FormatText && resolveColor(s.mode, s.w),
		})
	}
	return &Logger{c: c, name: cfg.name, fields: cfg.fields}
}

// Nop returns a logger that discards everything.
func Nop() *Logger { return New(WithLevel(LevelOff), WithOutput(io.Discard, ColorNever)) }

var (
	defMu  sync.RWMutex
	defLog *Logger
)

// Default returns the process-wide logger, creating it on first use.
func Default() *Logger {
	defMu.RLock()
	l := defLog
	defMu.RUnlock()
	if l != nil {
		return l
	}
	defMu.Lock()
	defer defMu.Unlock()
	if defLog == nil {
		defLog = New()
	}
	return defLog
}

// SetDefault replaces the process-wide logger.
func SetDefault(l *Logger) {
	defMu.Lock()
	defLog = l
	defMu.Unlock()
}

// With returns a child logger carrying extra permanent fields.
func (l *Logger) With(fields ...Field) *Logger {
	nf := make([]Field, 0, len(l.fields)+len(fields))
	nf = append(nf, l.fields...)
	nf = append(nf, fields...)
	return &Logger{c: l.c, name: l.name, fields: nf}
}

// Named returns a child logger with a (dot-joined) component name.
func (l *Logger) Named(name string) *Logger {
	if l.name != "" && name != "" {
		name = l.name + "." + name
	} else if name == "" {
		name = l.name
	}
	return &Logger{c: l.c, name: name, fields: l.fields}
}

// SetLevel changes the minimum level at runtime (affects all children).
func (l *Logger) SetLevel(lv Level) { l.c.level.Store(int32(lv)) }

// Level returns the current minimum level.
func (l *Logger) Level() Level { return Level(l.c.level.Load()) }

// Enabled reports whether a record of the given level would be emitted.
func (l *Logger) Enabled(lv Level) bool { return lv >= l.Level() && lv < LevelOff }

// Close flushes and closes every file sink.
func (l *Logger) Close() error {
	l.c.mu.Lock()
	defer l.c.mu.Unlock()
	var first error
	for _, s := range l.c.sinks {
		if s.closer != nil {
			if err := s.closer.Close(); err != nil && first == nil {
				first = err
			}
		}
	}
	return first
}

// Structured level methods.

func (l *Logger) Trace(msg string, f ...Field) { l.log(LevelTrace, msg, f) }
func (l *Logger) Debug(msg string, f ...Field) { l.log(LevelDebug, msg, f) }
func (l *Logger) Info(msg string, f ...Field)  { l.log(LevelInfo, msg, f) }
func (l *Logger) Warn(msg string, f ...Field)  { l.log(LevelWarn, msg, f) }
func (l *Logger) Error(msg string, f ...Field) { l.log(LevelError, msg, f) }

// Printf-style level methods.

func (l *Logger) Tracef(format string, a ...any) { l.logf(LevelTrace, format, a) }
func (l *Logger) Debugf(format string, a ...any) { l.logf(LevelDebug, format, a) }
func (l *Logger) Infof(format string, a ...any)  { l.logf(LevelInfo, format, a) }
func (l *Logger) Warnf(format string, a ...any)  { l.logf(LevelWarn, format, a) }
func (l *Logger) Errorf(format string, a ...any) { l.logf(LevelError, format, a) }

// Printf makes *Logger usable wherever a standard-library style logger
// is expected. The record is emitted at Info level.
func (l *Logger) Printf(format string, a ...any) { l.logf(LevelInfo, format, a) }

func (l *Logger) logf(lv Level, format string, a []any) {
	if !l.Enabled(lv) {
		return
	}
	l.emit(lv, fmt.Sprintf(format, a...), nil, 3)
}

func (l *Logger) log(lv Level, msg string, f []Field) {
	if !l.Enabled(lv) {
		return
	}
	l.emit(lv, msg, f, 3)
}

func (l *Logger) emit(lv Level, msg string, extra []Field, skip int) {
	e := Entry{Time: time.Now(), Level: lv, Name: l.name, Message: msg}
	if n := len(l.fields) + len(extra); n > 0 {
		e.Fields = make([]Field, 0, n)
		e.Fields = append(e.Fields, l.fields...)
		e.Fields = append(e.Fields, extra...)
	}
	if l.c.caller {
		if _, file, line, ok := runtime.Caller(skip); ok {
			e.Caller = filepath.Base(file) + ":" + itoa(line)
		}
	}
	bp := l.c.pool.Get().(*[]byte)
	defer l.c.pool.Put(bp)

	l.c.mu.Lock()
	defer l.c.mu.Unlock()
	for _, s := range l.c.sinks {
		if lv < s.min {
			continue
		}
		buf := (*bp)[:0]
		if s.format == FormatJSON {
			buf = appendJSON(buf, e)
		} else {
			buf = appendText(buf, e, l.c.timeFmt, s.color)
		}
		_, _ = s.w.Write(buf)
		*bp = buf
	}
	for _, h := range l.c.hooks {
		h(e)
	}
}

// Writer returns an io.Writer that logs every written line at lv.
func (l *Logger) Writer(lv Level) io.Writer { return &lineWriter{l: l, lv: lv} }

// StdLogger adapts the logger to the standard library type, for APIs
// that insist on *log.Logger.
func (l *Logger) StdLogger(lv Level) *log.Logger { return log.New(l.Writer(lv), "", 0) }

type lineWriter struct {
	l  *Logger
	lv Level
}

func (w *lineWriter) Write(p []byte) (int, error) {
	if w.l.Enabled(w.lv) {
		for _, line := range strings.Split(strings.TrimRight(string(p), "\r\n"), "\n") {
			w.l.emit(w.lv, line, nil, 4)
		}
	}
	return len(p), nil
}

// Package-level helpers using the default logger.

func Trace(msg string, f ...Field) { Default().log(LevelTrace, msg, f) }
func Debug(msg string, f ...Field) { Default().log(LevelDebug, msg, f) }
func Info(msg string, f ...Field)  { Default().log(LevelInfo, msg, f) }
func Warn(msg string, f ...Field)  { Default().log(LevelWarn, msg, f) }
func Error(msg string, f ...Field) { Default().log(LevelError, msg, f) }

func Infof(format string, a ...any)  { Default().logf(LevelInfo, format, a) }
func Warnf(format string, a ...any)  { Default().logf(LevelWarn, format, a) }
func Errorf(format string, a ...any) { Default().logf(LevelError, format, a) }
func Debugf(format string, a ...any) { Default().logf(LevelDebug, format, a) }
