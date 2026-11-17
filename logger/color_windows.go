//go:build windows

package logger

import (
	"os"
	"syscall"
	"unsafe"
)

const enableVirtualTerminalProcessing = 0x0004

var (
	kernel32    = syscall.NewLazyDLL("kernel32.dll")
	procGetMode = kernel32.NewProc("GetConsoleMode")
	procSetMode = kernel32.NewProc("SetConsoleMode")
)

// enableVT switches the console into ANSI processing mode.
func enableVT(f *os.File) bool {
	h := f.Fd()
	var mode uint32
	if r, _, _ := procGetMode.Call(h, uintptr(unsafe.Pointer(&mode))); r == 0 {
		return false
	}
	if mode&enableVirtualTerminalProcessing != 0 {
		return true
	}
	r, _, _ := procSetMode.Call(h, uintptr(mode|enableVirtualTerminalProcessing))
	return r != 0
}
