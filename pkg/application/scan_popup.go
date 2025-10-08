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

type ScanPopup struct {
	// app: used for moveToMainWindow
	app *Application

	query           string // use only in doMainQueryNoArg and lookup
	icon            *qt.QIcon
	popup           *qt.QWidget
	mode            dictmgr.SearchMode
	dragRelativePos *qt.QPoint
	font            *qt.QFont
	headerLabel     *HeaderLabel
	articleView     *ArticleView
}

func (p *ScanPopup) init() {
	popup := p.popup
	popup.SetWindowFlag(qt.FramelessWindowHint | qt.WindowStaysOnTopHint | qt.Tool)
	popup.SetWindowIcon(p.icon)

	p.headerLabel = NewHeaderLabel(p.doQuery)
	p.headerLabel.SetFont(p.font)
	p.headerLabel.SetMouseTracking(true)

	p.articleView = NewArticleView(p.doQuery)
	p.articleView.SetFont(p.font)
	p.articleView.SetupCustomHandlers()

	popup.MoveWithQPoint(qt.QCursor_Pos())
	popup.OnKeyPressEvent(p.onPopupKeyPress)

	for _, widget := range []HasOnMouseEvents{
		p.headerLabel,
		popup,
	} {
		widget.OnMousePressEvent(p.onDragMousePress)
		widget.OnMouseMoveEvent(p.onDragMouseMove)
		widget.OnMouseReleaseEvent(p.onDragMouseRelease)
	}

	popupLayout := qt.NewQVBoxLayout(popup)

	headerBox := qt.NewQHBoxLayout2()

	closeButton := qt.NewQPushButton3("Close")
	closeButton.SetFont(p.font)
	closeButton.OnClicked(func() {
		_ = popup.Close()
	})
	mainButton := qt.NewQPushButton3("Main")
	mainButton.SetFont(p.font)
	mainButton.OnClicked(p.moveToMainWindow)

	// favoriteButton := NewFavoriteButton(app.favoriteButtonClicked)
	// favoriteButton.SetToolTips(
	// 	"Add this term to favorites",
	// 	"Remove this term from favorites",
	// )

	headerButtonBox := qt.NewQVBoxLayout2()
	headerButtonBox.AddWidget(closeButton.QWidget)
	headerButtonBox.AddWidget(mainButton.QWidget)
	headerButtonBox.AddStretch()

	headerBox.AddWidget2(p.headerLabel.QWidget, 10)
	headerBox.AddStretch()
	headerBox.AddLayout(headerButtonBox.QLayout)

	popupLayout.AddLayout(headerBox.QLayout)
	// headerBox.AddWidget(favoriteButton.QWidget)
	popupLayout.AddWidget2(p.articleView.QWidget, 10)
}

func (p *ScanPopup) doQuery(query string) {
	results := dictmgr.LookupHTML(query, conf, p.mode, resultFlags, 0)
	if len(results) == 0 {
		slog.Info("scan popup", "min_score", conf.ScanPopupMinScore, "score", 0, "query", query)
		if conf.ScanPopupMinScore > 0 {
			return
		}
		p.articleView.SetHtml(fmt.Sprintf("No results for %#v", query))
		p.popup.SetWindowTitle(query)
		return
	}
	res := results[0]
	slog.Info("scan popup", "min_score", conf.ScanPopupMinScore, "score", res.Score()/2, "query", query)
	if conf.ScanPopupMinScore > int(res.Score())/2 {
		return
	}
	p.articleView.SetResult(res)
	p.headerLabel.SetResult(res)
	p.popup.SetWindowTitle(res.Terms()[0])
	// favoriteButton.SetChecked(app.favoritesWidget.HasFavorite(res.Terms()[0]))
	p.autoResize()
}

func (p *ScanPopup) moveToMainWindow() {
	p.popup.Close()
	p.app.window.Show()
	p.app.window.ActivateWindow()
	p.app.doQuery(p.query)
}

func (p *ScanPopup) onPopupKeyPress(super func(*qt.QKeyEvent), event *qt.QKeyEvent) {
	if event.Key() == escape {
		p.popup.Close()
		return
	}
	super(event)
}

func (p *ScanPopup) onDragMousePress(super func(*qt.QMouseEvent), event *qt.QMouseEvent) {
	if event.Button() != qt.LeftButton {
		super(event)
		return
	}
	p.dragRelativePos = event.Pos()
}

func (p *ScanPopup) onDragMouseMove(super func(*qt.QMouseEvent), event *qt.QMouseEvent) {
	if p.dragRelativePos == nil {
		super(event)
		return
	}
	p.popup.Move(
		event.GlobalX()-p.dragRelativePos.X(),
		event.GlobalY()-p.dragRelativePos.Y(),
	)
}

func (p *ScanPopup) onDragMouseRelease(super func(*qt.QMouseEvent), event *qt.QMouseEvent) {
	p.dragRelativePos = nil
	super(event)
}

func (p *ScanPopup) autoResize() {
	p.popup.Resize(conf.ScanPopupWidth, conf.ScanPopupHeight)
}

func (app *Application) onScanPopupCloseEvent(super func(*qt.QCloseEvent), event *qt.QCloseEvent) {
	app.scanPopupCount.Add(-1)
}

func NewScanPopup(
	app *Application,
	icon *qt.QIcon,
	query string,
	mode dictmgr.SearchMode,
) *ScanPopup {
	font := configFontWithFactor(conf.ScanPopupFontSizeFactor)
	popup := qt.NewQWidget2()
	popup.OnCloseEvent(app.onScanPopupCloseEvent)
	popup.SetFont(font)
	p := &ScanPopup{
		app:   app,
		query: query,
		icon:  icon,
		popup: popup,
		mode:  mode,
		font:  font,
	}
	p.init()

	return p
}

func (p *ScanPopup) Run() {
	p.doQuery(p.query)
	p.popup.Show()
}

func (app *Application) scanPopup(query string) {
	if conf.ScanPopupMaxCount > 0 && app.scanPopupCount.Load() >= conf.ScanPopupMaxCount {
		return
	}
	app.scanPopupCount.Add(1)

	query = strings.TrimSpace(query)
	query = strings.Trim(query, punctuation)
	if query == "" {
		return
	}
	mode, valid := dictmgr.SearchModeByName(conf.ScanPopupMode)
	if !valid {
		slog.Error("invalid scan_popup_mode", "value", conf.ScanPopupMode)
	}

	p := NewScanPopup(app, app.icon, query, mode)
	p.Run()
}
