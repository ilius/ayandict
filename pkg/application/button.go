package application

import (
	"fmt"

	"github.com/therecipe/qt/widgets"
)

func NewPNGIconTextButton(label string, imageName string) *widgets.QPushButton {
	icon, err := loadPNGIcon(imageName)
	if err != nil {
		fmt.Println(err)
	}
	if icon == nil {
		return widgets.NewQPushButton2(label, nil)
	}
	return widgets.NewQPushButton3(icon, label, nil)
}
