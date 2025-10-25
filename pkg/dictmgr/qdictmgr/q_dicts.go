package qdictmgr

import (
	"github.com/ilius/ayandict/v3/pkg/config"
	"github.com/ilius/ayandict/v3/pkg/dictmgr"
	"github.com/ilius/ayandict/v3/pkg/dictmgr/internal/dicts"
	"github.com/ilius/ayandict/v3/pkg/mysoundex"
	qt "github.com/mappu/miqt/qt6"
)

func loadingDictsPopup(conf *config.Config) *qt.QLabel {
	popup := qt.NewQLabel6(
		`<span style="font-size:xx-large;">Loading dictionaries</span>`,
		nil,
		qt.SplashScreen,
	)
	// Qt__SplashScreen makes it centered on screen
	popup.SetFrameStyle(int(qt.QFrame__Raised | qt.QFrame__Shadow(qt.QFrame__Panel)))
	popup.SetAlignment(qt.AlignCenter)
	popup.SetFixedSize2(300, 100)
	popup.SetWindowModality(qt.WindowModal)
	popup.SetStyleSheet(conf.LoadingPopupStyleStr)
	popup.Show()
	qt.QCoreApplication_ProcessEvents()
	return popup
}

func InitDicts(conf *config.Config, popup bool) {
	if popup {
		popupLabel := loadingDictsPopup(conf)
		defer popupLabel.Delete()
	}
	if conf.SoundexWordsFile != "" {
		ss := mysoundex.NewSoundexSearcher(conf.SoundexWordsFile)
		dictmgr.SetSoundexSearcher(ss)
		go ss.Load()
	}
	dicts.InitDicts(conf)
}
