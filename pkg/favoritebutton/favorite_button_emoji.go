package favoritebutton

import (
	"fmt"
	"log/slog"

	"github.com/ilius/ayandict/v3/pkg/config"
	qt "github.com/mappu/miqt/qt6"
)

func NewTextFavoriteButton(
	conf *config.Config,
	onClick func(bool),
	emoji bool,
) *TextFavoriteButton {
	activeHue := conf.FavoriteButtonHue
	qButton := qt.NewQPushButton2()
	if emoji {
		qButton.OnResizeEvent(func(super func(event *qt.QResizeEvent), event *qt.QResizeEvent) {
			qButton.SetFixedWidth(event.Size().Height())
		})
		qButton.SetText("★") // ★ ⭑ ✯ ☆ ⭐
	} else {
		qButton.SetText(" fav ")
	}
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

func hsvString(h int, s int, v int) string {
	color := qt.NewQColor()
	color.SetHsv(h, s, v)
	return color.ToQVariant().ToString()
}

// returns activeColor, inactiveColor
func getEmojiButtonColors(button *qt.QWidget, activeHue int) (string, string) {
	bgColor := button.Palette().Color(
		qt.QPalette__Normal,
		qt.QPalette__Base,
	)
	bgValue := bgColor.Value()
	slog.Debug(
		"EmojiFavoriteButton",
		"bgValue", bgValue,
		"bgColor", bgColor.ToQVariant().ToString(),
	)
	if bgValue < 128 { // dark theme
		return hsvString(
				activeHue,
				100,
				255,
			), hsvString(
				0,
				0,
				240,
			)
	}
	// light theme
	return hsvString(
			activeHue,
			100,
			bgValue*3/4,
		), hsvString(
			0,
			0,
			bgValue/2,
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

func (b *TextFavoriteButton) Button() *qt.QPushButton {
	return b.QPushButton
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
