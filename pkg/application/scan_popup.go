package application

import (
	"fmt"
	"log/slog"

	"github.com/ilius/ayandict/v3/pkg/dictmgr"
	"github.com/ilius/ayandict/v3/pkg/qplatform"
	commons "github.com/ilius/go-dict-commons"
	qt "github.com/mappu/miqt/qt6"
)

type QCloseEventFunc = func(super func(*qt.QCloseEvent), event *qt.QCloseEvent)

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
	articleView *ArticleView

	// set by doQuery
	results []commons.SearchResultIface

	// set by doQuery, gotoNextResult, gotoPrevResult
	resultIndex int

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
	popup.SetAttribute(qt.WA_DeleteOnClose)
	popup.SetWindowIcon(p.icon)
	popup.SetFont(font)

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

	popup.OnMousePressEvent(p.onMousePress)
	popup.OnMouseMoveEvent(p.onMouseMove)
	popup.OnMouseReleaseEvent(p.onMouseRelease)

	p.articleView.Browser.OnMousePressEvent(func(super func(ev *qt.QMouseEvent), ev *qt.QMouseEvent) {
		popup.ActivateWindow()
		super(ev)
	})

	popupLayout := qt.NewQVBoxLayout(popup)

	closeButton := qt.NewQPushButton3(" close ")
	closeButton.SetFont(font)
	closeButton.OnClicked(func() {
		_ = popup.Close()
	})
	mainButton := qt.NewQPushButton3(" main ")
	mainButton.SetFont(font)
	mainButton.OnClicked(p.moveToMainWindow)

	nextButton := qt.NewQPushButton3(" next ")
	nextButton.SetFont(font)
	nextButton.OnClicked(p.gotoNextResult)

	prevButton := qt.NewQPushButton3(" prev ")
	prevButton.SetFont(font)
	prevButton.OnClicked(p.gotoPrevResult)

	// favoriteButton := qt.NewQPushButton3("favorite")
	// favoriteButton.SetFont(font)
	// favoriteButton.SetCheckable(true)
	// favoriteButton.OnToggled(func(checked bool) {
	// 	// p.favoritesWidget.SetFavorite(term, checked)
	// })

	// favoriteButton := NewFavoriteButton(app.favoriteButtonClicked)
	// favoriteButton.SetToolTips(
	// 	"Add this term to favorites",
	// 	"Remove this term from favorites",
	// )

	headerBox := qt.NewQHBoxLayout2()
	{
		label := qt.NewQLabel2()
		label.SetCursor(qt.NewQCursor2(qt.DragMoveCursor))
		headerBox.AddWidget2(label.QWidget, 1)
	}
	headerBox.AddWidget(nextButton.QWidget)
	headerBox.AddWidget(prevButton.QWidget)
	headerBox.AddWidget(mainButton.QWidget)
	headerBox.AddWidget(closeButton.QWidget)

	popupLayout.SetContentsMargins(5, 0, 5, 5)
	popupLayout.SetSpacing(0)
	headerBox.SetContentsMargins(0, 0, 0, 0)
	headerBox.SetSpacing(0)

	// buttons have no ContentsMargins by default
	// setting "margin: 0px;" stylesheet mimimizes both the height and width
	// margin: top, right, bottom, left
	// const smallButtonSS = "margin: 0px 5px 0px 5px;"
	const smallButtonSS = "margin: 0px;"
	closeButton.SetStyleSheet(smallButtonSS)
	mainButton.SetStyleSheet(smallButtonSS)
	nextButton.SetStyleSheet(smallButtonSS)
	prevButton.SetStyleSheet(smallButtonSS)
	// closeButton.SetContentsMargins(5, 0, 5, 0)
	// mainButton.SetContentsMargins(5, 0, 5, 0)

	popupLayout.AddLayout(headerBox.QLayout)
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
	p.results = results
	p.resultIndex = 0
	res := results[0]
	slog.Info("scan popup", "min_score", conf.ScanPopupMinScore, "score", res.Score()/2, "query", query)
	if conf.ScanPopupMinScore > int(res.Score())/2 {
		return
	}
	p.articleView.SetResultWithHeader(res)
	p.popup.SetWindowTitle(res.Terms()[0])
	// favoriteButton.SetChecked(app.favoritesWidget.HasFavorite(res.Terms()[0]))
	p.autoResize()
	if conf.ScanPopupHistory {
		p.addHistory(query)
	}
}

func (p *ScanPopup) gotoNextResult() {
	index := p.resultIndex + 1
	if index > len(p.results)-1 {
		return
	}
	res := p.results[index]
	p.resultIndex = index
	p.articleView.SetResultWithHeader(res)
	// if conf.ScanPopupHistory {
	// 	p.addHistory(res.Terms()[0])
	// }
}

func (p *ScanPopup) gotoPrevResult() {
	index := p.resultIndex - 1
	if index < 0 {
		return
	}
	res := p.results[index]
	p.resultIndex = index
	p.articleView.SetResultWithHeader(res)
	// if conf.ScanPopupHistory {
	// 	p.addHistory(res.Terms()[0])
	// }
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

func (p *ScanPopup) onMousePress(super func(*qt.QMouseEvent), event *qt.QMouseEvent) {
	if event.Button() != qt.LeftButton {
		super(event)
		return
	}
	p.dragPos = event.Pos()
	p.popup.ActivateWindow()
}

func (p *ScanPopup) onMouseMove(super func(*qt.QMouseEvent), event *qt.QMouseEvent) {
	if p.dragPos == nil {
		super(event)
		return
	}
	p.popup.Move(
		event.GlobalX()-p.dragPos.X(),
		event.GlobalY()-p.dragPos.Y(),
	)
}

func (p *ScanPopup) onMouseRelease(super func(*qt.QMouseEvent), event *qt.QMouseEvent) {
	p.dragPos = nil
	super(event)
}

func (p *ScanPopup) autoResize() {
	p.popup.Resize(conf.ScanPopupWidth, conf.ScanPopupHeight)
}
