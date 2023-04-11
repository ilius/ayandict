package qerr

import (
	"fmt"
	"log"

	"github.com/therecipe/qt/widgets"
)

var ShowQtError = false

func Error(args ...any) {
	msg := fmt.Sprint(args...)
	log.Println(msg) // TODO: stderr
	if !ShowQtError {
		return
	}
	defer func() {
		r := recover()
		if r != nil {
			log.Println(r)
		}
	}()
	d := widgets.NewQErrorMessage(nil)
	d.ShowMessage(msg)
}

func Errorf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	log.Println(msg) // TODO: stderr
	if !ShowQtError {
		return
	}
	defer func() {
		r := recover()
		if r != nil {
			log.Println(r)
		}
	}()
	d := widgets.NewQErrorMessage(nil)
	d.ShowMessage(msg)
}
