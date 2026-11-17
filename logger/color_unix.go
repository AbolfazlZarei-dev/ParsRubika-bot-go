//go:build !windows

package logger

import "os"

// enableVT is a no-op on platforms whose terminals understand ANSI.
func enableVT(*os.File) bool { return true }
