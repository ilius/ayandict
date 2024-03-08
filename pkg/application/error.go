package application

import (
	"log"

	"github.com/ilius/qt/widgets"
)

func showErrorMessage(msg string) {
	defer func() {
		r := recover()
		if r != nil {
			log.Println(r)
		}
	}()
	d := widgets.NewQErrorMessage(nil)
	d.ShowMessage(msg)
}
