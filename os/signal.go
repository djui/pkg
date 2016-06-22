package os

import (
	"io"
	"os"
	"os/signal"
	"syscall"
)

// WaitForSignals listens looping for syscall signals until the given signal
// handler returns true. Common uses are CTRL-T, CTRL-C, CTRL-D, CTRL-Z.
//
// Note SIGQUIT (CTRL-\) is not handle on purpose to not accidentally override
// Go's stacktracing.
func WaitForSignals(handler func(os.Signal) bool) {
	done := make(chan struct{}, 1)
	sigc := make(chan os.Signal, 5)

	signal.Notify(sigc,
		syscall.SIGUSR1, // syscall.SIGINFO, // BSD only
		syscall.SIGHUP, syscall.SIGUSR2,
		syscall.SIGINT, syscall.SIGTERM,
	)
	for sig := range sigc {
		switch sig {
		case // syscall.SIGINFO, // Ctrl-t
			syscall.SIGUSR1,
			syscall.SIGHUP,
			syscall.SIGUSR2,
			syscall.SIGINT, // Ctrl-c
			syscall.SIGTERM:
			signal.Stop(sigc)
			if handler(sig) {
				done <- struct{}{}
				return
			}
		}
	}

	go func() {
		for {
			select {
			case <-done:
				return
			default:
				if _, err := os.Stdin.Read([]byte{0}); err == io.EOF { // Ctrl-d
					sigc <- syscall.SIGHUP
				}
			}
		}
	}()
}
