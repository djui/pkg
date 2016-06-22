package exec

import (
	"os/exec"
	"syscall"
)

// ExitCode returns the exit code of a given error or `false, 0` if the given
// error is no ExitError.
func ExitCode(err error) (int, bool) {
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				return status.ExitStatus(), true
			}
		}
	}

	return 0, false
}

// IsExitError returns if a given error is an ExitError.
func IsExitError(err error) bool {
	if err == nil {
		return false
	}

	_, ok := err.(*exec.ExitError)
	return ok
}
