package application

import (
	"fmt"

	qt "github.com/mappu/miqt/qt6"
)

func NewPNGIconTextButton(label string, imageName string) *qt.QPushButton {
	icon, err := loadPNGIcon(imageName)
	if err != nil {
		fmt.Println(err)
	}
	if icon == nil {
		return qt.NewQPushButton3(label)
	}
	return qt.NewQPushButton4(icon, label)
}
