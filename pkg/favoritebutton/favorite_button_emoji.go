package favoritebutton

import (
	"fmt"
	"log/slog"

	"github.com/ilius/ayandict/v3/pkg/config"
	"github.com/ilius/ayandict/v3/pkg/qtutils"
	qt "github.com/mappu/miqt/qt6"
)

func NewColoredEmojiFavoriteButton(
	conf *config.Config,
	onClick func(bool),
) *TextFavoriteButton {
	activeHue := conf.FavoriteButtonHue
	qButton := qt.NewQPushButton2()
	qButton.OnResizeEvent(func(super func(event *qt.QResizeEvent), event *qt.QResizeEvent) {
		qButton.SetFixedWidth(event.Size().Height())
	})
	qButton.SetText("★") // ★ ⭑ ✯ ☆ ⭐
	button := &TextFavoriteButton{
		QPushButton: qButton,
	}
	activeColor, inactiveColor := getEmojiButtonColors(qButton.QWidget, activeHue)
	button.activeStyleSheet = fmt.Sprintf(
		"QPushButton{color: %s; margin: 0px;}",
		activeColor,
	)
	button.inactiveStyleSheet = fmt.Sprintf(
		"QPushButton{color: %s; margin: 0px;}",
		inactiveColor,
	)
	button.OnClicked(func() {
		onClick(button.ToggleChecked())
	})
	return button
}

// returns activeColor, inactiveColor
func getEmojiButtonColors(button *qt.QWidget, activeHue int) (string, string) {
	bg := button.Palette().Color(
		qt.QPalette__Normal,
		qt.QPalette__Base,
	)
	fg := button.Palette().Color(
		qt.QPalette__Normal,
		qt.QPalette__ButtonText,
	)
	slog.Debug(
		"EmojiFavoriteButton",
		"bg.V", bg.Value(),
		"fg.V", fg.Value(),
		"bg.L", bg.Lightness(),
		"fg.L", fg.Lightness(),
		"bg", bg.ToQVariant().ToString(),
		"fg", fg.ToQVariant().ToString(),
	)
	// "dusk" color scheme: bg.V=127 fg.V=0 bg.L=127 fg.L=0 bg=#7f7f7f fg=#000000
	if bg.Value() < 130 { // dark BG
		return qtutils.HSVString(
				activeHue,
				110,
				255-bg.Value()/6,
			), qtutils.HSVString(
				0,
				0,
				fg.Value()*94/100,
			)
	}
	// light BG
	return qtutils.HSVString(
			activeHue,
			100,
			bg.Value()*3/4,
		), qtutils.HSVString(
			0,
			0,
			bg.Value()/2,
		)
}

type TextFavoriteButton struct {
	*qt.QPushButton
	checked bool

	activeStyleSheet   string
	inactiveStyleSheet string

	inactiveTooltip string
	activeTooltip   string
}

func (b *TextFavoriteButton) QWidget() *qt.QWidget {
	return b.QPushButton.QWidget
}

func (b *TextFavoriteButton) SetChecked(checked bool) {
	b.checked = checked
	if checked {
		b.SetToolTip(b.activeTooltip)
		b.SetStyleSheet(b.activeStyleSheet)
	} else {
		b.SetToolTip(b.inactiveTooltip)
		b.SetStyleSheet(b.inactiveStyleSheet)
	}
}

func (b *TextFavoriteButton) ToggleChecked() bool {
	b.SetChecked(!b.checked)
	return b.checked
}

func (b *TextFavoriteButton) SetToolTips(inactive string, active string) {
	b.inactiveTooltip = inactive
	b.activeTooltip = active
	b.SetToolTip(inactive)
}
