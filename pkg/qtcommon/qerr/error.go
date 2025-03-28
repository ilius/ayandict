package qerr

import (
	"fmt"
)

// ShowMessage: set in GUI application
var ShowMessage = func(_ string) {}

func Error(args ...any) {
	msg := fmt.Sprint(args...)
	ShowMessage(msg)
}

func Errorf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	ShowMessage(msg)
}
