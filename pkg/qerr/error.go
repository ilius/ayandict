package qerr

import (
	"fmt"
	"log/slog"
)

// ShowMessage: set in GUI application
var ShowMessage = func(msg string) {}

func Error(args ...any) {
	msg := fmt.Sprint(args...)
	slog.Error(msg)
	ShowMessage(msg)
}

func Errorf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	slog.Error(msg)
	ShowMessage(msg)
}
