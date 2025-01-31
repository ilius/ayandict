package application

import (
	"log/slog"

	qt "github.com/mappu/miqt/qt6"
)

func showErrorMessage(msg string) {
	defer func() {
		r := recover()
		if r != nil {
			slog.Error("Panic", "r", r)
		}
	}()
	d := qt.NewQErrorMessage(nil)
	d.ShowMessage(msg)
}
