package logger

import (
	"io"
	"os"
	"strings"
)

// ColorMode controls when ANSI colours are emitted.
type ColorMode int

const (
	// ColorAuto enables colours when the output is a capable terminal.
	ColorAuto ColorMode = iota
	// ColorAlways forces colours on.
	ColorAlways
	// ColorNever disables colours.
	ColorNever
)

// ANSI escape sequences used by the text formatter.
const (
	ansiReset     = "\x1b[0m"
	ansiBold      = "\x1b[1m"
	ansiDim       = "\x1b[2m"
	ansiRed       = "\x1b[31m"
	ansiGreen     = "\x1b[32m"
	ansiYellow    = "\x1b[33m"
	ansiBlue      = "\x1b[34m"
	ansiMagenta   = "\x1b[35m"
	ansiCyan      = "\x1b[36m"
	ansiGray      = "\x1b[90m"
	ansiBrightRed = "\x1b[91m"
	ansiWhiteOnRd = "\x1b[1;97;41m"
	ansiBlackOnGr = "\x1b[30;42m"
	ansiBlackOnBl = "\x1b[30;44m"
	ansiBlackOnYl = "\x1b[30;43m"
	ansiBlackOnGy = "\x1b[97;100m"
)

// levelStyle returns the badge and message colours of a level.
func levelStyle(l Level) (badge, text string) {
	switch l {
	case LevelTrace:
		return ansiBlackOnGy, ansiGray
	case LevelDebug:
		return ansiBlackOnBl, ""
	case LevelInfo:
		return ansiBlackOnGr, ""
	case LevelWarn:
		return ansiBlackOnYl, ansiYellow
	default:
		return ansiWhiteOnRd, ansiBrightRed
	}
}

// resolveColor decides whether w should receive colours.
func resolveColor(mode ColorMode, w io.Writer) bool {
	switch mode {
	case ColorAlways:
		return true
	case ColorNever:
		return false
	}
	if _, ok := os.LookupEnv("NO_COLOR"); ok {
		return false
	}
	if v := os.Getenv("FORCE_COLOR"); v != "" && v != "0" {
		return true
	}
	if strings.EqualFold(os.Getenv("TERM"), "dumb") {
		return false
	}
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	return isTerminal(f) && enableVT(f)
}

// isTerminal reports whether f is attached to a character device.
func isTerminal(f *os.File) bool {
	st, err := f.Stat()
	if err != nil {
		return false
	}
	return st.Mode()&os.ModeCharDevice != 0
}
