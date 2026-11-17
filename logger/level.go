package logger

import (
	"fmt"
	"strings"
)

// Level is the severity of a log record.
type Level int32

// Severity levels, from the most verbose to the most critical.
const (
	LevelTrace Level = iota - 1
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelOff
)

// String returns the upper-case name of the level.
func (l Level) String() string {
	switch l {
	case LevelTrace:
		return "TRACE"
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelOff:
		return "OFF"
	}
	return fmt.Sprintf("LEVEL(%d)", int32(l))
}

// short returns the fixed-width (5 column) label used in text output.
func (l Level) short() string {
	switch l {
	case LevelTrace:
		return "TRACE"
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO "
	case LevelWarn:
		return "WARN "
	case LevelError:
		return "ERROR"
	case LevelOff:
		return "OFF  "
	}
	return "     "
}

// ParseLevel converts a textual level ("debug", "warn", ...) to a Level.
func ParseLevel(s string) (Level, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "trace":
		return LevelTrace, nil
	case "debug":
		return LevelDebug, nil
	case "", "info":
		return LevelInfo, nil
	case "warn", "warning":
		return LevelWarn, nil
	case "error", "err":
		return LevelError, nil
	case "off", "none", "silent":
		return LevelOff, nil
	}
	return LevelInfo, fmt.Errorf("logger: unknown level %q", s)
}
