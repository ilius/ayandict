package application

import "github.com/ilius/qt/widgets"

type FavoriteButton struct {
	*widgets.QPushButton
}

func NewFavoriteButton() *FavoriteButton {
	button := NewPNGIconTextButton("", "favorite.png")
	button.SetCheckable(true)
	return &FavoriteButton{
		QPushButton: button,
	}
}
