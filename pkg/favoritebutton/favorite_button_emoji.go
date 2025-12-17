package favoritebutton

import (
	"fmt"
	"log/slog"

	"github.com/ilius/ayandict/v3/pkg/config"
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
	qButton.SetStyleSheet("QPushButton{margin: 0px;}")
	activeColor, inactiveColor := getEmojiButtonColors(qButton.QWidget, activeHue)
	button := &TextFavoriteButton{
		QPushButton:   qButton,
		activeColor:   activeColor,
		inactiveColor: inactiveColor,
	}
	if !config.PrivateMode {
		button.OnClicked(func() {
			onClick(button.ToggleChecked())
		})
	}
	return button
}

// returns activeColor, inactiveColor
func getEmojiButtonColors(button *qt.QWidget, activeHue int) (*qt.QColor, *qt.QColor) {
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
	active := qt.NewQColor()
	inactive := qt.NewQColor()
	if bg.Value() < 130 { // dark BG
		active.SetHsv(
			activeHue,
			110,
			255-bg.Value()/6,
		)
		inactive.SetHsv(
			0,
			0,
			fg.Value()*94/100,
		)
		return active, inactive
	}
	// light BG
	active.SetHsv(
		activeHue,
		100,
		bg.Value()*3/4,
	)
	inactive.SetHsv(
		0,
		0,
		bg.Value()/2,
	)
	return active, inactive
}

type TextFavoriteButton struct {
	*qt.QPushButton
	checked bool

	activeColor   *qt.QColor
	inactiveColor *qt.QColor

	inactiveTooltip string
	activeTooltip   string
}

func (b *TextFavoriteButton) QWidget() *qt.QWidget {
	return b.QPushButton.QWidget
}

func (b *TextFavoriteButton) SetChecked(checked bool) {
	b.checked = checked
	var color *qt.QColor
	if checked {
		b.SetToolTip(b.activeTooltip)
		color = b.activeColor
	} else {
		b.SetToolTip(b.inactiveTooltip)
		color = b.inactiveColor
	}
	b.SetStyleSheet(fmt.Sprintf(
		"QPushButton{color: %s; margin: 0px;}",
		color.ToQVariant().ToString(),
	))

	// this also works (instead of SetStyleSheet), but still probably not reliable
	// b.Palette().SetColor(qt.QPalette__All, qt.QPalette__ButtonText, color)
	// b.Palette().SetColor(qt.QPalette__All, qt.QPalette__Text, color)
	// b.Update()
	// b.Repaint()
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
