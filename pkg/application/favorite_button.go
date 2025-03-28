package application

import (
	"log/slog"

	qt "github.com/mappu/miqt/qt6"
)

type FavoriteButton struct {
	*qt.QPushButton
	checked      bool
	activeIcon   *qt.QIcon
	inactiveIcon *qt.QIcon

	inactiveTooltip string
	activeTooltip   string
}

func (b *FavoriteButton) SetChecked(checked bool) {
	b.checked = checked
	if checked {
		b.SetIcon(b.activeIcon)
		b.SetToolTip(b.activeTooltip)
	} else {
		b.SetIcon(b.inactiveIcon)
		b.SetToolTip(b.inactiveTooltip)
	}
}

func (b *FavoriteButton) ToggleChecked() {
	b.SetChecked(!b.checked)
}

func (b *FavoriteButton) SetToolTips(inactive string, active string) {
	b.inactiveTooltip = inactive
	b.activeTooltip = active
	b.SetToolTip(inactive)
}

func NewFavoriteButton(onClick func(bool)) *FavoriteButton {
	activeIcon, err := loadPNGIcon("favorite-active-64.png")
	if err != nil {
		slog.Error("error loading icon favorite-active-64.png: " + err.Error())
		panic(err)
	}
	inactiveIcon, err := loadPNGIcon("favorite-64.png")
	if err != nil {
		slog.Error("error loading icon favorite-64.png: " + err.Error())
		panic(err)
	}
	qButton := qt.NewQPushButton4(inactiveIcon, "")
	qButton.OnResizeEvent(func(super func(event *qt.QResizeEvent), event *qt.QResizeEvent) {
		iconSize := event.Size().Height() * 2 / 3
		qButton.SetIconSize(qt.NewQSize2(iconSize, iconSize))
	})
	button := &FavoriteButton{
		QPushButton:  qButton,
		activeIcon:   activeIcon,
		inactiveIcon: inactiveIcon,
	}
	button.OnClicked(func() {
		button.ToggleChecked()
		onClick(button.checked)
	})
	return button
}
