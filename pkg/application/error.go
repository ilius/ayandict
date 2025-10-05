package application

import (
	"log"
	"os"

	qt "github.com/mappu/miqt/qt6"
)

var stdLogger = log.New(os.Stderr, "", log.LstdFlags)

func showErrorMessage(msg string) {
	defer func() {
		r := recover()
		if r != nil {
			stdLogger.Printf("Panic: %v", r)
		}
	}()
	d := qt.NewQErrorMessage(nil)
	d.ShowMessage(msg)
}
