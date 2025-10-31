package application

import (
	"bytes"
	"fmt"
	"html/template"
	"log/slog"

	"github.com/ilius/ayandict/v3/pkg/dictmgr"
	"github.com/ilius/ayandict/v3/pkg/qplatform"
	common "github.com/ilius/go-dict-commons"
	commons "github.com/ilius/go-dict-commons"
	qt "github.com/mappu/miqt/qt6"
)

type QCloseEventFunc = func(super func(*qt.QCloseEvent), event *qt.QCloseEvent)

var aboutToDragCursor = qt.PointingHandCursor

type ScanPopupAppInterface interface {
	OnScanPopupShow()
	ShowWindowAndQuery(query string)
	AddHistoryAndFrequency(query string)
	ShowAbout()
}

func NewScanPopup(
	query string,
	mode dictmgr.SearchMode,
	onCloseEvent QCloseEventFunc,
	appIface ScanPopupAppInterface,
) *ScanPopup {
	popup := qt.NewQWidget2()
	popup.OnCloseEvent(onCloseEvent)

	p := &ScanPopup{
		query: query,
		mode:  mode,
		popup: popup,
		app:   appIface,
	}
	p.init()
	p.popup.OnShowEvent(func(super func(*qt.QShowEvent), e *qt.QShowEvent) {
		appIface.OnScanPopupShow()
		super(e)
	})
	return p
}

type ScanPopup struct {
	// set by factory func:
	query string             // used in doMainQueryNoArg and Run
	mode  dictmgr.SearchMode // used in doQuery
	app   ScanPopupAppInterface

	popup *qt.QWidget

	// set by init method:
	headerTemplate *template.Template
	articleView    *ArticleView
	nextButton     *qt.QPushButton
	prevButton     *qt.QPushButton

	// set by doQuery
	results []commons.SearchResultIface

	// set by doQuery, gotoNextResult, gotoPrevResult
	resultIndex int

	// dragPos: position of mouse relative to window, when drag starts
	dragPos *qt.QPoint

	// header label that you use to drag
	headerLabel *qt.QLabel
}

func (p *ScanPopup) Run(pos *qt.QPoint, icon *qt.QIcon) {
	if pos == nil {
		screen := qt.QGuiApplication_Screens()[0]
		ss := screen.Size()
		p.popup.Move(
			(ss.Width()-conf.ScanPopupWidth)/2,
			(ss.Height()-conf.ScanPopupHeight)/2,
		)
	} else {
		screen := qt.QGuiApplication_ScreenAt(pos)
		p.popup.Move(pos.X(), pos.Y()+int(fontPixelSize(systemFont, screen)))
	}
	p.popup.SetWindowIcon(icon)

	p.doQuery(p.query)
}

func (p *ScanPopup) init() {
	font := configFontWithFactor(conf.ScanPopupFontSizeFactor)

	popup := p.popup
	flags := qt.WindowStaysOnTopHint | qt.Tool | qt.FramelessWindowHint
	if qplatform.CanMoveWindow() {
		if conf.ScanPopupBypassWindowManager {
			flags |= qt.BypassWindowManagerHint
		}
	}
	popup.SetWindowFlags(flags)
	popup.SetAttribute(qt.WA_DeleteOnClose)
	popup.SetFont(font)

	headerTemplate := template.New("popupheader")
	headerTemplate, err := headerTemplate.Parse(conf.ScanPopupHeaderTemplate)
	if err != nil {
		slog.Error("error parsing scan popup template", "err", err)
	}
	p.headerTemplate = headerTemplate

	p.articleView = NewArticleView(p.doQuery)
	p.articleView.Widget.SetFont(font)
	p.articleView.OnKeyPressEvent(p.onArticleViewKeyPressEvent)

	popup.OnKeyPressEvent(p.onKeyPress)

	popup.OnMousePressEvent(p.onMousePress)
	popup.OnMouseMoveEvent(p.onMouseMove)
	popup.OnMouseReleaseEvent(p.onMouseRelease)

	p.articleView.OnMousePressEvent(func(super func(ev *qt.QMouseEvent), ev *qt.QMouseEvent) {
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
	p.nextButton = nextButton

	prevButton := qt.NewQPushButton3(" prev ")
	prevButton.SetFont(font)
	prevButton.OnClicked(p.gotoPrevResult)
	p.prevButton = prevButton

	if conf.ScanPopupHeaderIcons {
		closeButton.SetText(" ❌ ")
		mainButton.SetText(" 📖 ")
		nextButton.SetText("  ↓  ")
		prevButton.SetText("  ↑  ")
	}

	closeButton.SetToolTip("Close (Esc)")
	mainButton.SetToolTip("Open in main window (Enter)")
	nextButton.SetToolTip("Next result (Alt+Down)")
	prevButton.SetToolTip("Previous result (Alt+Up)")

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
	headerBox.AddSpacing(10)
	{
		label := qt.NewQLabel2()
		label.SetAutoFillBackground(true)
		label.SetBackgroundRole(qt.QPalette__ToolTipBase)
		label.OnMousePressEvent(p.onMousePress)
		label.OnMouseMoveEvent(p.onMouseMove)
		label.OnMouseReleaseEvent(p.onMouseRelease)
		// label.SetFont(font)
		p.headerLabel = label
		label.SetCursor(qt.NewQCursor2(aboutToDragCursor))
		headerBox.AddWidget2(label.QWidget, 1)
	}
	headerBox.AddWidget(nextButton.QWidget)
	headerBox.AddWidget(prevButton.QWidget)
	headerBox.AddWidget(mainButton.QWidget)
	headerBox.AddWidget(closeButton.QWidget)

	nextButton.SetFocusPolicy(qt.NoFocus)
	prevButton.SetFocusPolicy(qt.NoFocus)
	mainButton.SetFocusPolicy(qt.NoFocus)
	closeButton.SetFocusPolicy(qt.NoFocus)

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

	popupLayout.AddWidget(p.outlineHeaderLayout(headerBox.QLayout))
	popupLayout.AddWidget2(p.articleView.Widget, 10)
}

func (p *ScanPopup) outlineHeaderLayout(layout *qt.QLayout) *qt.QWidget {
	frame := qt.NewQFrame2()
	frame.SetFrameShape(qt.QFrame__Box)
	frame.SetFrameShadow(qt.QFrame__Sunken)
	frame.SetContentsMargins(1, 1, 1, 1)
	frameLayout := qt.NewQHBoxLayout(frame.QWidget)
	frameLayout.AddLayout2(layout, 1)
	frameLayout.SetContentsMargins(0, 0, 0, 0)
	return frame.QWidget
}

func (p *ScanPopup) onQueryNoResult(message string, args ...any) {
	p.articleView.SetHtml(fmt.Sprintf(message, args...))
	p.nextButton.Hide()
	p.prevButton.Hide()
	p.popup.Show()
	p.popup.ActivateWindow()
}

type ScanPopupHeaderTemplateInput struct {
	DictName string
	Score    uint8
}

func (p *ScanPopup) setResult(res common.SearchResultIface) *qt.QTimer {
	if p.headerTemplate == nil {
		slog.Error("popup header template is not set")
		return nil
	}
	headerBuf := bytes.NewBuffer(nil)
	dictName := res.DictName()
	err := p.headerTemplate.Execute(headerBuf, ScanPopupHeaderTemplateInput{
		DictName: dictName,
		Score:    res.Score() >> 1,
	})
	if err != nil {
		slog.Error("error encoding popup header", "err", err)
		return nil
	}
	p.headerLabel.SetText(headerBuf.String())
	return p.articleView.SetPopupResult(
		res,
		conf.ScanPopupTermsStyle,
	)
}

func (p *ScanPopup) doQuery(query string) {
	p.query = query
	p.popup.SetWindowTitle(query)

	results := dictmgr.LookupHTML(query, conf, p.mode, resultFlags, 0)
	if len(results) == 0 {
		slog.Info("scan popup", "min_score", conf.ScanPopupMinScore, "score", 0, "query", query)
		p.onQueryNoResult("No results for %#v", query)
		return
	}
	p.results = results
	p.resultIndex = 0
	res := results[0]
	slog.Info("scan popup", "min_score", conf.ScanPopupMinScore, "score", res.Score()/2, "query", query)
	if conf.ScanPopupMinScore > int(res.Score())/2 {
		p.onQueryNoResult(
			"Top result for %#v has score of %%%v",
			query, res.Score()/2,
		)
		return
	}
	p.nextButton.Show()
	p.prevButton.Show()
	playTimer := p.setResult(res)
	// favoriteButton.SetChecked(app.favoritesWidget.HasFavorite(res.Terms()[0]))
	p.autoResize()
	if conf.ScanPopupHistory {
		p.app.AddHistoryAndFrequency(query)
	}
	p.popup.Show()
	p.popup.ActivateWindow()
	if playTimer != nil {
		qt.QCoreApplication_ProcessEvents()
		playTimer.Start(0)
	}
}

func (p *ScanPopup) gotoNextResult() {
	index := p.resultIndex + 1
	if index > len(p.results)-1 {
		return
	}
	res := p.results[index]
	p.resultIndex = index
	playTimer := p.setResult(res)
	// if conf.ScanPopupHistory {
	// 	p.addHistory(res.Terms()[0])
	// }
	if playTimer != nil {
		playTimer.Start(0)
	}
}

func (p *ScanPopup) gotoPrevResult() {
	index := p.resultIndex - 1
	if index < 0 {
		return
	}
	res := p.results[index]
	p.resultIndex = index
	playTimer := p.setResult(res)
	// if conf.ScanPopupHistory {
	// 	p.addHistory(res.Terms()[0])
	// }
	if playTimer != nil {
		playTimer.Start(0)
	}
}

func (p *ScanPopup) moveToMainWindow() {
	p.popup.Close()
	p.app.ShowWindowAndQuery(p.query)
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

func (p *ScanPopup) onArticleViewKeyPressEvent(super func(event *qt.QKeyEvent), event *qt.QKeyEvent) {
	switch event.Key() {
	case int(qt.Key_Up):
		if event.Modifiers()&qt.AltModifier > 0 {
			p.gotoPrevResult()
		} else {
			super(event)
		}
	case int(qt.Key_Down):
		if event.Modifiers()&qt.AltModifier > 0 {
			p.gotoNextResult()
		} else {
			super(event)
		}
	case int(qt.Key_F1):
		p.popup.Close()
		p.app.ShowAbout()
	default:
		super(event)
	}
}

func (p *ScanPopup) onMousePress(super func(*qt.QMouseEvent), event *qt.QMouseEvent) {
	if event.Button() != qt.LeftButton {
		super(event)
		return
	}
	p.popup.ActivateWindow()
	if qplatform.CanMoveWindow() {
		p.dragPos = event.WindowPos().ToPoint()
		p.headerLabel.SetCursor(qt.NewQCursor2(qt.DragMoveCursor))
	} else {
		p.popup.WindowHandle().StartSystemMove()
	}
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
	p.headerLabel.SetCursor(qt.NewQCursor2(aboutToDragCursor))
	super(event)
}

func (p *ScanPopup) autoResize() {
	p.popup.Resize(conf.ScanPopupWidth, conf.ScanPopupHeight)
}
