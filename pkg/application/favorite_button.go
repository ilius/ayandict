package application

import (
	"github.com/ilius/ayandict/v2/pkg/qtcommon/qerr"
	"github.com/ilius/qt/widgets"
)

type FavoriteButton struct {
	*widgets.QPushButton
}

func NewFavoriteButton() *FavoriteButton {
	icon, err := loadPNGIcon("favorite.png")
	if err != nil {
		qerr.Error(err)
		panic(err)
	}
	if icon == nil {
		qerr.Error("error loading favorite.png icon: icon is nil")
		panic("error loading favorite.png icon: icon is nil")
	}
	button := widgets.NewQPushButton3(icon, "", nil)
	button.SetCheckable(true)
	return &FavoriteButton{
		QPushButton: button,
	}
}
