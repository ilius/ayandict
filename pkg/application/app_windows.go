//go:build windows

package application

import (
	"os"
	"os/signal"
	"syscall"
)

var (
	modkernel32               = syscall.NewLazyDLL("kernel32.dll")
	procSetConsoleCtrlHandler = modkernel32.NewProc("SetConsoleCtrlHandler")
)

const (
	CTRL_C_EVENT        = 0
	CTRL_CLOSE_EVENT    = 2
	CTRL_LOGOFF_EVENT   = 5
	CTRL_SHUTDOWN_EVENT = 6
)

func init() {
	// handle Ctrl+C
	signal.Notify(exitSignalChan, os.Interrupt)

	// Optional: Handle console close -> "virtual SIGTERM"
	h := syscall.NewCallback(func(ctrl uint) uintptr {
		switch ctrl {
		case CTRL_C_EVENT:
			exitSignalChan <- os.Interrupt
			return 1
		case CTRL_CLOSE_EVENT, CTRL_LOGOFF_EVENT, CTRL_SHUTDOWN_EVENT:
			exitSignalChan <- syscall.SIGTERM
			return 1
		}
		return 0
	})
	_, _, _ = procSetConsoleCtrlHandler.Call(h, 1)
}
