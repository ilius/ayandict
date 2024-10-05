package application

import (
	"log/slog"

	"github.com/ilius/qt/widgets"
)

func showErrorMessage(msg string) {
	defer func() {
		r := recover()
		if r != nil {
			slog.Error("Panic", "r", r)
		}
	}()
	d := widgets.NewQErrorMessage(nil)
	d.ShowMessage(msg)
}
