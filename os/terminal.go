// +build darwin dragonfly freebsd linux,!appengine netbsd openbsd

package os

import (
	"io"
	"os"
	"syscall"
	"unsafe"
)

// IsStdinPiped checks if data is being piped to os.Stdin or if os.Stdin is from
// a terminal.
func IsStdinPiped() bool {
	stat, _ := os.Stdin.Stat()
	return stat.Mode()&os.ModeCharDevice == 0
}

// IsStdinReadable checks if os.Stdin is open for reading in a blocking way.
func IsStdinReadable() bool {
	// unix.SetNonblock(int(os.Stdin.Fd()), true)
	_, err := os.Stdin.Read([]byte{0})
	// unix.SetNonblock(int(os.Stdin.Fd()), false)
	return err != io.EOF
}

// IsStdinTerminal returns true if stdin is a terminal.
func IsStdinTerminal() bool {
	return IsTerminal(os.Stdin.Fd())
}

// IsTerminal returns true if the given file descriptor is a terminal.
func IsTerminal(fd uintptr) bool {
	var termios syscall.Termios
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, fd, syscall.TIOCGETA, uintptr(unsafe.Pointer(&termios)), 0, 0, 0)
	return err == 0
}
