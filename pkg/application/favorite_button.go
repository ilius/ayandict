package application

import (
	"github.com/ilius/ayandict/v2/pkg/qtcommon/qerr"
	"github.com/ilius/qt/gui"
	"github.com/ilius/qt/widgets"
)

type FavoriteButton struct {
	*widgets.QPushButton
	checked      bool
	activeIcon   *gui.QIcon
	inactiveIcon *gui.QIcon
}

func (b *FavoriteButton) SetChecked(checked bool) {
	b.checked = checked
	if checked {
		b.SetIcon(b.activeIcon)
	} else {
		b.SetIcon(b.inactiveIcon)
	}
}

func (b *FavoriteButton) ToggleChecked() {
	b.SetChecked(!b.checked)
}

func NewFavoriteButton(onClick func(bool)) *FavoriteButton {
	activeIcon, err := loadPNGIcon("favorite-active-64.png")
	if err != nil {
		qerr.Error(err)
		panic(err)
	}
	inactiveIcon, err := loadPNGIcon("favorite-64.png")
	if err != nil {
		qerr.Error(err)
		panic(err)
	}
	qButton := widgets.NewQPushButton3(inactiveIcon, "", nil)
	// FIXME: reduce internal padding / border width
	button := &FavoriteButton{
		QPushButton:  qButton,
		activeIcon:   activeIcon,
		inactiveIcon: inactiveIcon,
	}
	button.ConnectClicked(func(checked bool) {
		button.ToggleChecked()
		onClick(button.checked)
	})
	return button
}
