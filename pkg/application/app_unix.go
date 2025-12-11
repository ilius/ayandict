//go:build !windows

package application

import (
	"os/signal"
	"syscall"
)

func init() {
	signal.Notify(
		exitSignalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
}
