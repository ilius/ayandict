//go:build !darwin

package application

import (
	qt "github.com/mappu/miqt/qt6"
)

func showErrorMessage(msg string) {
	defer func() {
		r := recover()
		if r != nil {
			stdLogger.Printf("Panic: %v", r)
		}
	}()
	d := qt.NewQErrorMessage(nil)
	d.OnFinished(func(result int) {
		d.Destroy()
	})
	d.ShowMessage(msg)
	d.Exec()
}
