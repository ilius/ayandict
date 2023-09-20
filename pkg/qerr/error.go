package qerr

import (
	"fmt"
	"log"
)

// ShowMessage: set in GUI application
var ShowMessage = func(msg string) {}

func Error(args ...any) {
	msg := fmt.Sprint(args...)
	log.Println(msg) // TODO: stderr
	ShowMessage(msg)
}

func Errorf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	log.Println(msg) // TODO: stderr
	ShowMessage(msg)
}
