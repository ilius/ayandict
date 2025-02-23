package application

import (
	"log/slog"

	"github.com/ilius/qt/widgets"
)

func NewPNGIconTextButton(label string, imageName string) *widgets.QPushButton {
	icon, err := loadPNGIcon(imageName)
	if err != nil {
		slog.Error("error loading png icon", "imageName", imageName, "err", err)
		return widgets.NewQPushButton2(label, nil)
	}
	if icon == nil {
		slog.Error("error loading png icon: icon is nil", "imageName", imageName)
		return widgets.NewQPushButton2(label, nil)
	}
	return widgets.NewQPushButton3(icon, label, nil)
}
