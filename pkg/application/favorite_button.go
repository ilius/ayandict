package application

import (
	"github.com/ilius/ayandict/v2/pkg/qtcommon/qerr"
	"github.com/ilius/qt/core"
	"github.com/ilius/qt/gui"
	"github.com/ilius/qt/widgets"
)

type FavoriteButton struct {
	*widgets.QPushButton
	checked      bool
	activeIcon   *gui.QIcon
	inactiveIcon *gui.QIcon

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
		qerr.Error(err)
		panic(err)
	}
	inactiveIcon, err := loadPNGIcon("favorite-64.png")
	if err != nil {
		qerr.Error(err)
		panic(err)
	}
	qButton := widgets.NewQPushButton3(inactiveIcon, "", nil)
	qButton.ConnectResizeEvent(func(event *gui.QResizeEvent) {
		iconSize := event.Size().Height() * 4 / 5
		qButton.SetIconSize(core.NewQSize2(iconSize, iconSize))
	})
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
