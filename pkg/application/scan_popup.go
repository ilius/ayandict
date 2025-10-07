package application

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/ilius/ayandict/v3/pkg/dictmgr"
	qt "github.com/mappu/miqt/qt6"
)

type HasOnMouseEvents interface {
	OnMousePressEvent(func(super func(*qt.QMouseEvent), event *qt.QMouseEvent))
	OnMouseMoveEvent(func(super func(*qt.QMouseEvent), event *qt.QMouseEvent))
	OnMouseReleaseEvent(func(super func(*qt.QMouseEvent), event *qt.QMouseEvent))
	ObjectName() string
}

func (app *Application) setupScanPopup() {
	clipboard := qt.QGuiApplication_Clipboard()
	clipboard.OnChanged(func(mode qt.QClipboard__Mode) {
		if mode == qt.QClipboard__Clipboard && !conf.ScanPopupClipboard {
			return
		}
		if mode == qt.QClipboard__Selection && !conf.ScanPopupSelection {
			return
		}
		app.scanPopup(clipboard.TextWithMode(mode))
	})
}

func (app *Application) scanPopup(query string) {
	if conf.ScanPopupMaxCount > 0 && app.scanPopupCount.Load() >= conf.ScanPopupMaxCount {
		return
	}
	app.scanPopupCount.Add(1)
	popup := qt.NewQWidget2()
	popup.SetWindowFlag(qt.FramelessWindowHint | qt.WindowStaysOnTopHint)
	popup.SetWindowIcon(app.icon)
	popup.OnCloseEvent(func(super func(*qt.QCloseEvent), event *qt.QCloseEvent) {
		app.scanPopupCount.Add(-1)
	})

	query = strings.TrimSpace(query)
	query = strings.Trim(query, punctuation)
	if query == "" {
		return
	}
	mode, valid := dictmgr.SearchModeByName(conf.ScanPopupMode)
	if !valid {
		slog.Error("invalid scan_popup_mode", "value", conf.ScanPopupMode)
	}

	font := *ConfigFont()
	font.SetPixelSize(int(float64(font.PixelSize()) * conf.ScanPopupFontSizeFactor))
	popup.SetFont(&font)

	headerLabel := CreateHeaderLabel(app)
	headerLabel.SetFont(&font)

	articleView := NewArticleView(app)
	articleView.SetFont(&font)

	popupLayout := qt.NewQVBoxLayout(popup)

	headerBox := qt.NewQHBoxLayout2()

	queryInMainWindow := func() {
		popup.Close()
		app.window.Show()
		app.window.ActivateWindow()
		app.entry.SetText(query)
		onQuery(query, app.queryArgs, false)
	}
	headerLabel.doQuery = func(term string) { queryInMainWindow() }

	closeButton := qt.NewQPushButton3("Close")
	closeButton.SetFont(&font)
	closeButton.OnClicked(func() {
		popup.Close()
	})
	mainButton := qt.NewQPushButton3("Main")
	mainButton.SetFont(&font)
	mainButton.OnClicked(queryInMainWindow)

	// favoriteButton := NewFavoriteButton(app.favoriteButtonClicked)
	// favoriteButton.SetToolTips(
	// 	"Add this term to favorites",
	// 	"Remove this term from favorites",
	// )

	headerButtonBox := qt.NewQVBoxLayout2()
	headerButtonBox.AddWidget(closeButton.QWidget)
	headerButtonBox.AddWidget(mainButton.QWidget)
	headerButtonBox.AddStretch()

	headerBox.AddWidget2(headerLabel.QWidget, 10)
	headerBox.AddStretch()
	headerBox.AddLayout(headerButtonBox.QLayout)

	popupLayout.AddLayout(headerBox.QLayout)
	// headerBox.AddWidget(favoriteButton.QWidget)
	popupLayout.AddWidget2(articleView.QWidget, 10)

	popup.MoveWithQPoint(qt.QCursor_Pos())
	popup.OnKeyPressEvent(func(super func(*qt.QKeyEvent), event *qt.QKeyEvent) {
		if event.Key() == escape {
			popup.Close()
			return
		}
		super(event)
	})

	headerLabel.SetMouseTracking(true)

	var dragRelativePos *qt.QPoint
	for _, widget := range []HasOnMouseEvents{
		headerLabel,
		popup,
	} {
		widget.OnMousePressEvent(func(super func(*qt.QMouseEvent), event *qt.QMouseEvent) {
			if event.Button() != qt.LeftButton {
				super(event)
				return
			}
			dragRelativePos = event.Pos()
		})
		widget.OnMouseMoveEvent(func(super func(*qt.QMouseEvent), event *qt.QMouseEvent) {
			if dragRelativePos == nil {
				super(event)
				return
			}
			popup.Move(
				event.GlobalX()-dragRelativePos.X(),
				event.GlobalY()-dragRelativePos.Y(),
			)
		})
		widget.OnMouseReleaseEvent(func(super func(*qt.QMouseEvent), event *qt.QMouseEvent) {
			dragRelativePos = nil
			super(event)
		})
	}

	results := dictmgr.LookupHTML(query, conf, mode, resultFlags, 0)
	if len(results) == 0 {
		slog.Info("scan popup", "min_score", conf.ScanPopupMinScore, "score", 0, "query", query)
		if conf.ScanPopupMinScore > 0 {
			return
		}
		articleView.SetHtml(fmt.Sprintf("No results for %#v", query))
		popup.SetWindowTitle(query)
	} else {
		res := results[0]
		slog.Info("scan popup", "min_score", conf.ScanPopupMinScore, "score", res.Score()/2, "query", query)
		if conf.ScanPopupMinScore > int(res.Score())/2 {
			return
		}
		articleView.SetResult(res)
		headerLabel.SetResult(res)
		popup.SetWindowTitle(res.Terms()[0])
		// favoriteButton.SetChecked(app.favoritesWidget.HasFavorite(res.Terms()[0]))
	}

	popup.Resize(conf.ScanPopupWidth, conf.ScanPopupHeight)
	popup.Show()
}
