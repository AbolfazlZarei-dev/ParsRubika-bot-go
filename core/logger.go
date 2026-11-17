package core

// Logger is the minimal logging contract shared by every ParsRubika
// component. *log.Logger and *logger.Logger both satisfy it.
type Logger interface {
	Printf(format string, args ...any)
}

// LevelLogger is the optional richer contract. When the installed
// logger also implements it, components emit properly levelled
// records instead of plain Printf lines.
type LevelLogger interface {
	Logger
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)
}

// Level identifies the severity used by Logf.
type Level int

// Severity levels understood by Logf.
const (
	LogDebug Level = iota
	LogInfo
	LogWarn
	LogError
)

// Logf writes to l at the given level, degrading to Printf for loggers
// that only implement the basic contract. A nil logger is ignored.
func Logf(l Logger, lv Level, format string, args ...any) {
	if l == nil {
		return
	}
	if ll, ok := l.(LevelLogger); ok {
		switch lv {
		case LogDebug:
			ll.Debugf(format, args...)
		case LogInfo:
			ll.Infof(format, args...)
		case LogWarn:
			ll.Warnf(format, args...)
		default:
			ll.Errorf(format, args...)
		}
		return
	}
	l.Printf(format, args...)
}
