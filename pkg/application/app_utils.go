package application

import (
	"github.com/ilius/ayandict/v3/pkg/config"
	qt "github.com/mappu/miqt/qt6"
)

// returns basePx which is %66 of the font size in pixels,
// I'm using it for spacing between widgets
// kinda like "em" in html, but probably not exactly the same
func (app *Application) baseFontPixelSize() float32 {
	return float32(fontPixelSize(
		qt.QApplication_Font(),
		app.window.Screen(),
	) * 0.66)
}

func (app *Application) newIconTextButton(label string, pix qt.QStyle__StandardPixmap) *qt.QPushButton {
	return qt.NewQPushButton4(
		app.style.StandardIcon(
			pix,
			app.bottomBoxStyleOpt.QStyleOption,
			nil,
		),
		label,
	)
}

func (app *Application) makeAboutButton(conf *config.Config) *qt.QPushButton {
	aboutButtonLabel := "About"
	if conf.ReduceMinimumWindowWidth {
		aboutButtonLabel = "\u200c"
	}
	aboutButton := app.newIconTextButton(aboutButtonLabel, qt.QStyle__SP_MessageBoxInformation)
	aboutButton.OnClicked(app.ShowAbout)
	return aboutButton
}

// func (app *Application) tableWidgetItem(text string) *qt.QTableWidgetItem {
// 	item := qt.NewQTableWidgetItem2(text)
// 	item.SetTextAlignment(0)
// 	return item
// }
