package application

import (
	"log/slog"

	qt "github.com/mappu/miqt/qt6"
)

func NewPNGIconTextButton(label string, imageName string) *qt.QPushButton {
	icon, err := loadPNGIcon(imageName)
	if err != nil {
		slog.Error("error loading png icon", "imageName", imageName, "err", err)
		return qt.NewQPushButton3(label)
	}
	if icon == nil {
		slog.Error("error loading png icon: icon is nil", "imageName", imageName)
		return qt.NewQPushButton3(label)
	}
	return qt.NewQPushButton4(icon, label)
}
