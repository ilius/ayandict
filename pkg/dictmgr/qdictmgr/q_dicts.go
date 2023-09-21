package qdictmgr

import (
	"github.com/ilius/ayandict/v2/pkg/config"
	"github.com/ilius/ayandict/v2/pkg/dictmgr/internal/dicts"
	"github.com/ilius/qt/core"
	"github.com/ilius/qt/widgets"
)

func loadingDictsPopup(conf *config.Config) *widgets.QLabel {
	popup := widgets.NewQLabel2(
		`<span style="font-size:xx-large;">Loading dictionaries</span>`,
		nil,
		core.Qt__SplashScreen,
	)
	// Qt__SplashScreen makes it centered on screen
	popup.SetFrameStyle(int(widgets.QFrame__Raised | widgets.QFrame__Shadow(widgets.QFrame__Panel)))
	popup.SetAlignment(core.Qt__AlignCenter)
	popup.SetFixedSize2(300, 100)
	popup.SetWindowModality(core.Qt__WindowModal)
	popup.SetStyleSheet(conf.PopupStyleStr)
	popup.Show()
	core.QCoreApplication_ProcessEvents(core.QEventLoop__AllEvents)
	return popup
}

func InitDicts(conf *config.Config, popup bool) {
	if popup {
		popupLabel := loadingDictsPopup(conf)
		defer popupLabel.Destroy(true, true)
	}
	dicts.InitDicts(conf)
}
