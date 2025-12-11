package favoritebutton

import (
	"github.com/ilius/ayandict/v3/pkg/config"
	qt "github.com/mappu/miqt/qt6"
)

func NewMinimalFavoriteButton(
	_ *config.Config,
	onClick func(bool),
) *MinimalFavoriteButton {
	qButton := qt.NewQPushButton2()
	button := &MinimalFavoriteButton{
		QPushButton:  qButton,
		inactiveText: " fav",
		activeText:   "✔fav",
	}
	button.OnClicked(func() {
		onClick(button.ToggleChecked())
	})
	button.SetChecked(false)
	return button
}

type MinimalFavoriteButton struct {
	*qt.QPushButton
	checked bool

	inactiveText string
	activeText   string

	inactiveTooltip string
	activeTooltip   string
}

func (b *MinimalFavoriteButton) QWidget() *qt.QWidget {
	return b.QPushButton.QWidget
}

func (b *MinimalFavoriteButton) SetChecked(checked bool) {
	b.checked = checked
	if checked {
		b.SetToolTip(b.activeTooltip)
		b.SetText(b.activeText)
	} else {
		b.SetToolTip(b.inactiveTooltip)
		b.SetText(b.inactiveText)
	}
}

func (b *MinimalFavoriteButton) ToggleChecked() bool {
	b.SetChecked(!b.checked)
	return b.checked
}

func (b *MinimalFavoriteButton) SetToolTips(inactive string, active string) {
	b.inactiveTooltip = inactive
	b.activeTooltip = active
	b.SetToolTip(inactive)
}
