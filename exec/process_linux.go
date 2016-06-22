package exec

import (
	"os/exec"
	"syscall"
)

// SetPdeathsig set the parent-death signal for an executable command. The
// effect of this is that all child processes spawned by the parent process will
// exit when the parent process exits.
func SetPdeathsig(cmd *exec.Cmd, signal syscall.Signal) {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.Pdeathsig = signal
}
