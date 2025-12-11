package favoritebutton

import (
	"log/slog"

	"github.com/ilius/ayandict/v3/pkg/config"
	qt "github.com/mappu/miqt/qt6"
)

func NewImageFavoriteButton(
	_ *config.Config,
	onClick func(bool),
) *ImageFavoriteButton {
	activeIcon, err := loadPNGIcon("favorite-green-64.png")
	if err != nil {
		slog.Error("error loading icon favorite-green-64.png: " + err.Error())
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
	button := &ImageFavoriteButton{
		QPushButton:  qButton,
		activeIcon:   activeIcon,
		inactiveIcon: inactiveIcon,
	}
	button.OnClicked(func() {
		onClick(button.ToggleChecked())
	})
	return button
}

type ImageFavoriteButton struct {
	*qt.QPushButton
	checked      bool
	activeIcon   *qt.QIcon
	inactiveIcon *qt.QIcon

	inactiveTooltip string
	activeTooltip   string
}

func (b *ImageFavoriteButton) QWidget() *qt.QWidget {
	return b.QPushButton.QWidget
}

func (b *ImageFavoriteButton) SetChecked(checked bool) {
	b.checked = checked
	if checked {
		b.SetIcon(b.activeIcon)
		b.SetToolTip(b.activeTooltip)
	} else {
		b.SetIcon(b.inactiveIcon)
		b.SetToolTip(b.inactiveTooltip)
	}
}

func (b *ImageFavoriteButton) ToggleChecked() bool {
	b.SetChecked(!b.checked)
	return b.checked
}

func (b *ImageFavoriteButton) SetToolTips(inactive string, active string) {
	b.inactiveTooltip = inactive
	b.activeTooltip = active
	b.SetToolTip(inactive)
}
