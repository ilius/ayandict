package application

import (
	"fmt"
	"log/slog"

	"github.com/ilius/ayandict/v3/pkg/dictmgr"
	"github.com/ilius/ayandict/v3/pkg/qplatform"
	qt "github.com/mappu/miqt/qt6"
)

type QCloseEventFunc = func(super func(*qt.QCloseEvent), event *qt.QCloseEvent)

type HasOnMouseEvents interface {
	OnMousePressEvent(func(super func(*qt.QMouseEvent), event *qt.QMouseEvent))
	OnMouseMoveEvent(func(super func(*qt.QMouseEvent), event *qt.QMouseEvent))
	OnMouseReleaseEvent(func(super func(*qt.QMouseEvent), event *qt.QMouseEvent))
	ObjectName() string
}

func NewScanPopup(
	query string,
	mode dictmgr.SearchMode,
	pos *qt.QPoint,
	icon *qt.QIcon,
	showInMain func(query string),
	addHistory func(query string),
	onCloseEvent QCloseEventFunc,
) *ScanPopup {
	popup := qt.NewQWidget2()
	popup.OnCloseEvent(onCloseEvent)

	p := &ScanPopup{
		query:      query,
		mode:       mode,
		pos:        pos,
		icon:       icon,
		popup:      popup,
		showInMain: showInMain,
		addHistory: addHistory,
	}
	p.init()
	return p
}

type ScanPopup struct {
	// set by factory func:
	query      string             // used in doMainQueryNoArg and Run
	mode       dictmgr.SearchMode // used in doQuery
	pos        *qt.QPoint
	icon       *qt.QIcon
	popup      *qt.QWidget
	showInMain func(query string)
	addHistory func(query string)

	// set by init method:
	headerLabel *HeaderLabel
	articleView *ArticleView

	// dragPos: position of mouse relative to window, when drag starts
	dragPos *qt.QPoint
}

func (p *ScanPopup) Run() {
	p.doQuery(p.query)
	p.popup.Show()
	p.popup.ActivateWindow()
}

func (p *ScanPopup) init() {
	font := configFontWithFactor(conf.ScanPopupFontSizeFactor)

	popup := p.popup
	flags := qt.WindowStaysOnTopHint | qt.Tool
	if qplatform.CanMoveWindow() {
		flags |= qt.FramelessWindowHint
		if conf.ScanPopupBypassWindowManager {
			flags |= qt.BypassWindowManagerHint
		}
	} else {
		slog.Info("ScanPopup: normal (bordered), platform does not support moving window")
	}
	popup.SetWindowFlags(flags)
	popup.SetWindowIcon(p.icon)
	popup.SetFont(font)

	p.headerLabel = NewHeaderLabel(p.doQuery)
	p.headerLabel.SetFont(font)
	p.headerLabel.SetMouseTracking(true)

	p.articleView = NewArticleView(p.doQuery)
	p.articleView.Widget.SetFont(font)
	p.articleView.SetupCustomHandlers()

	pos := p.pos
	if pos == nil {
		screen := qt.QGuiApplication_Screens()[0]
		ss := screen.Size()
		popup.Move(
			(ss.Width()-conf.ScanPopupWidth)/2,
			(ss.Height()-conf.ScanPopupHeight)/2,
		)
	} else {
		screen := qt.QGuiApplication_ScreenAt(pos)
		popup.Move(pos.X(), pos.Y()+int(fontPixelSize(systemFont, screen)))
	}

	popup.OnKeyPressEvent(p.onKeyPress)
	p.headerLabel.OnKeyPressEvent(p.onKeyPress)
	p.articleView.Browser.OnKeyPressEvent(p.onKeyPress)

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
	closeButton.SetFont(font)
	closeButton.OnClicked(func() {
		_ = popup.Close()
	})
	mainButton := qt.NewQPushButton3("Main")
	mainButton.SetFont(font)
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
	popupLayout.AddWidget2(p.articleView.Widget, 10)
}

func (p *ScanPopup) doQuery(query string) {
	p.query = query
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
	if conf.ScanPopupHistory {
		p.addHistory(query)
	}
}

func (p *ScanPopup) moveToMainWindow() {
	p.popup.Close()
	p.showInMain(p.query)
}

func (p *ScanPopup) onKeyPress(super func(*qt.QKeyEvent), event *qt.QKeyEvent) {
	switch event.Key() {
	case escape:
		p.popup.Close()
	case int(qt.Key_Return), int(qt.Key_Enter):
		if p.articleView.Searching() {
			super(event)
		} else {
			p.moveToMainWindow()
		}
	default:
		super(event)
	}
}

func (p *ScanPopup) onDragMousePress(super func(*qt.QMouseEvent), event *qt.QMouseEvent) {
	if event.Button() != qt.LeftButton {
		super(event)
		return
	}
	p.dragPos = event.Pos()
}

func (p *ScanPopup) onDragMouseMove(super func(*qt.QMouseEvent), event *qt.QMouseEvent) {
	if p.dragPos == nil {
		super(event)
		return
	}
	p.popup.Move(
		event.GlobalX()-p.dragPos.X(),
		event.GlobalY()-p.dragPos.Y(),
	)
}

func (p *ScanPopup) onDragMouseRelease(super func(*qt.QMouseEvent), event *qt.QMouseEvent) {
	p.dragPos = nil
	super(event)
}

func (p *ScanPopup) autoResize() {
	p.popup.Resize(conf.ScanPopupWidth, conf.ScanPopupHeight)
}
