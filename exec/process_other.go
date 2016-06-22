// +build !linux

package exec

import (
	"os/exec"
	"syscall"
)

// SetPdeathsig set the parent-death signal for an executable command. The
// effect of this is that all child processes spawned by the parent process will
// exit when the parent process exits.
//
// Currently a No-Op as exec.Cmd.SysProcAttr.Pdeathsig is not implemented on
// Darwin.
func SetPdeathsig(cmd *exec.Cmd, signal syscall.Signal) {
	// No-Op
}
