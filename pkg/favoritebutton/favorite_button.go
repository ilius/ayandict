package favoritebutton

import (
	"log/slog"

	"github.com/ilius/ayandict/v3/pkg/resources"
	"github.com/ilius/ayandict/v3/pkg/resourceutil"
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

func (b *FavoriteButton) Button() *qt.QPushButton {
	return b.QPushButton
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

func loadPNGIcon(filename string) (*qt.QIcon, error) {
	return resourceutil.LoadPNGIcon(resources.Res, filename)
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
		h := event.Size().Height()
		iconSize := h * 2 / 3
		qButton.SetIconSize(qt.NewQSize2(iconSize, iconSize))
		qButton.SetFixedWidth(h)
	})
	qButton.SetStyleSheet("margin: 0px;")
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
