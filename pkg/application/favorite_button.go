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
	button := widgets.NewQPushButton3(icon, "", nil)
	button.SetCheckable(true)
	return &FavoriteButton{
		QPushButton: button,
	}
}
