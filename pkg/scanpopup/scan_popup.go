package scanpopup

import (
	"bytes"
	"fmt"
	"html/template"
	"log/slog"

	common "codeberg.org/ilius/go-dict-commons"
	commons "codeberg.org/ilius/go-dict-commons"
	"github.com/ilius/ayandict/v3/pkg/articleview"
	"github.com/ilius/ayandict/v3/pkg/audiocache"
	"github.com/ilius/ayandict/v3/pkg/config"
	"github.com/ilius/ayandict/v3/pkg/dictmgr"
	"github.com/ilius/ayandict/v3/pkg/favoritebutton"
	"github.com/ilius/ayandict/v3/pkg/qplatform"
	"github.com/ilius/ayandict/v3/pkg/qtcommon"
	"github.com/ilius/ayandict/v3/pkg/qtutils"
	qt "github.com/mappu/miqt/qt6"
)

const resultFlags = uint32(
	common.ResultFlag_FixAudio |
		common.ResultFlag_FixFileSrc |
		common.ResultFlag_FixWordLink |
		common.ResultFlag_ColorMapping)

type QCloseEventFunc = func(super func(*qt.QCloseEvent), event *qt.QCloseEvent)

var aboutToDragCursor = qt.PointingHandCursor

type FavoriteButtonInterface interface {
	SetChecked(bool)
	SetToolTips(string, string)
	QWidget() *qt.QWidget
}

type ScanPopupAppInterface interface {
	OnScanPopupShow()
	OnScanPopupClose(super func(*qt.QCloseEvent), event *qt.QCloseEvent)
	ShowWindowAndQuery(query string)
	AddHistoryAndFrequency(query string)
	HasFavorite(term string) bool
	SetFavoriteFromPopup(term string, favorite bool)
	ShowAbout()
	AudioCache() *audiocache.AudioCache
}

func NewScanPopup(
	conf *config.Config,
	query string,
	mode dictmgr.SearchMode,
	onCloseEvent QCloseEventFunc,
	appIface ScanPopupAppInterface,
) *ScanPopup {
	popup := qt.NewQWidget2()
	popup.OnCloseEvent(onCloseEvent)

	p := &ScanPopup{
		conf:  conf,
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

func (p *ScanPopup) QueryPopup(query string) {
	p2 := NewScanPopup(
		p.conf,
		query,
		p.mode,
		p.app.OnScanPopupClose,
		p.app,
	)
	p2.Run(qt.QCursor_Pos(), p.popup.WindowIcon())
}

type ScanPopup struct {
	// set by factory func:
	conf  *config.Config
	query string             // used in doMainQueryNoArg and Run
	mode  dictmgr.SearchMode // used in doQuery
	app   ScanPopupAppInterface

	popup *qt.QWidget

	// set by init method:
	headerTemplate *template.Template
	articleView    *articleview.ArticleView
	nextButton     *qt.QPushButton
	prevButton     *qt.QPushButton
	favoriteButton FavoriteButtonInterface

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
			(ss.Width()-p.conf.ScanPopupWidth)/2,
			(ss.Height()-p.conf.ScanPopupHeight)/2,
		)
	} else {
		systemFont := qt.QApplication_Font()
		screen := qt.QGuiApplication_ScreenAt(pos)
		p.popup.Move(pos.X(), pos.Y()+int(qtutils.FontPixelSize(systemFont, screen)))
	}
	p.popup.SetWindowIcon(icon)

	p.Query(p.query)
}

func (app *ScanPopup) IsPopup() bool {
	return true
}

func (p *ScanPopup) configFontWithFactor(factor float64) *qt.QFont {
	font := *qtcommon.ConfigFont(p.conf)
	font.SetPixelSize(int(float64(font.PixelSize()) * factor))
	return &font
}

func (p *ScanPopup) newHeaderButton(text string) *qt.QPushButton {
	button := qt.NewQPushButton3(text)
	// buttons have no ContentsMargins by default
	// setting "margin: 0px;" stylesheet mimimizes both the height and width
	// margin: top, right, bottom, left
	// "margin: 0px 5px 0px 5px;"
	button.SetStyleSheet("margin: 0px;")
	if p.conf.ScanPopupHeaderIcons {
		button.OnResizeEvent(func(super func(*qt.QResizeEvent), event *qt.QResizeEvent) {
			super(event)
			button.SetFixedWidth(event.Size().Height())
		})
	}
	return button
}

func (p *ScanPopup) init() {
	font := p.configFontWithFactor(p.conf.ScanPopupFontSizeFactor)

	popup := p.popup
	flags := qt.WindowStaysOnTopHint | qt.Tool | qt.FramelessWindowHint
	if qplatform.CanMoveWindow() {
		if p.conf.ScanPopupBypassWindowManager {
			flags |= qt.BypassWindowManagerHint
		}
	}
	popup.SetWindowFlags(flags)
	popup.SetAttribute(qt.WA_DeleteOnClose)
	popup.SetFont(font)

	headerTemplate := template.New("popupheader")
	headerTemplate, err := headerTemplate.Parse(p.conf.ScanPopupHeaderTemplate)
	if err != nil {
		slog.Error("error parsing scan popup template", "err", err)
	}
	p.headerTemplate = headerTemplate

	p.articleView = articleview.NewArticleView(p.conf, p)
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

	nextButton := p.newHeaderButton(" next ")
	nextButton.OnClicked(p.gotoNextResult)
	nextButton.SetToolTip("Next result (Alt+Down)")
	p.nextButton = nextButton

	prevButton := p.newHeaderButton(" prev ")
	prevButton.OnClicked(p.gotoPrevResult)
	prevButton.SetToolTip("Previous result (Alt+Up)")
	p.prevButton = prevButton

	if p.conf.ScanPopupHeaderIcons {
		p.favoriteButton = favoritebutton.NewColoredEmojiFavoriteButton(
			p.conf,
			p.favoriteButtonClicked,
		)
	} else {
		p.favoriteButton = favoritebutton.NewMinimalFavoriteButton(
			p.conf,
			p.favoriteButtonClicked,
		)
	}
	p.favoriteButton.SetToolTips(
		"Add this term to favorites",
		"Remove this term from favorites",
	)

	mainButton := p.newHeaderButton(" main ")
	mainButton.OnClicked(p.moveToMainWindow)
	mainButton.SetToolTip("Open in main window (Enter)")

	closeButton := p.newHeaderButton(" close ")
	closeButton.OnClicked(func() {
		_ = popup.Close()
	})
	closeButton.SetToolTip("Close (Esc)")

	if p.conf.ScanPopupHeaderIcons {
		nextButton.SetText("▾") // ▾▼↓⬇⮯🔽⬇️
		prevButton.SetText("▴") // ▴▲↑⬆⮭🔼⬆️
		mainButton.SetText("📖") // 📖📔📗↩️🔍
		closeButton.SetText("❌")
		bg := p.popup.Palette().Color(qt.QPalette__Normal, qt.QPalette__Base)
		fg := p.popup.Palette().Color(qt.QPalette__Normal, qt.QPalette__ButtonText)
		c := qt.NewQColor()
		if bg.Value() < 130 { // dark BG
			c.SetHsv(0, 0, fg.Value()*92/100)
		} else {
			c.SetHsv(0, 0, bg.Value()/2)
		}
		nextButton.Palette().SetColor(qt.QPalette__All, qt.QPalette__ButtonText, c)
		prevButton.Palette().SetColor(qt.QPalette__All, qt.QPalette__ButtonText, c)
	}

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

	for _, button := range []*qt.QWidget{
		nextButton.QWidget,
		prevButton.QWidget,
		p.favoriteButton.QWidget(),
		mainButton.QWidget,
		closeButton.QWidget,
	} {
		button.SetSizePolicy2(qt.QSizePolicy__Minimum, qt.QSizePolicy__Expanding)
		button.SetFont(font)
		button.SetFocusPolicy(qt.NoFocus)
		headerBox.AddWidget(button)
	}

	popupLayout.SetContentsMargins(5, 0, 5, 5)
	popupLayout.SetSpacing(0)
	headerBox.SetContentsMargins(0, 0, 0, 0)
	headerBox.SetSpacing(0)

	popupLayout.AddWidget(p.outlineHeaderLayout(headerBox.QLayout))
	popupLayout.AddWidget2(p.articleView.Widget, 10)
}

func (p *ScanPopup) AudioCache() *audiocache.AudioCache {
	return p.app.AudioCache()
}

func (p *ScanPopup) favoriteButtonClicked(checked bool) {
	p.app.SetFavoriteFromPopup(
		p.results[p.resultIndex].Terms()[0],
		checked,
	)
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
	p.favoriteButton.QWidget().Hide()
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
	p.favoriteButton.SetChecked(p.app.HasFavorite(res.Terms()[0]))
	return p.articleView.SetPopupResult(
		res,
		p.conf.ScanPopupTermsStyle,
	)
}

func (p *ScanPopup) Query(query string) {
	p.query = query
	p.popup.SetWindowTitle(query)

	results := dictmgr.LookupHTML(query, p.conf, p.mode, resultFlags, 0)
	if len(results) == 0 {
		slog.Info("scan popup", "min_score", p.conf.ScanPopupMinScore, "score", 0, "query", query)
		p.onQueryNoResult("No results for %#v", query)
		return
	}
	p.results = results
	p.resultIndex = 0
	res := results[0]
	slog.Info("scan popup", "min_score", p.conf.ScanPopupMinScore, "score", res.Score()/2, "query", query)
	if p.conf.ScanPopupMinScore > int(res.Score())/2 {
		p.onQueryNoResult(
			"Top result for %#v has score of %%%v",
			query, res.Score()/2,
		)
		return
	}
	p.nextButton.Show()
	p.prevButton.Show()
	p.favoriteButton.QWidget().Show()
	playTimer := p.setResult(res)
	p.autoResize()
	if p.conf.ScanPopupHistory {
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
	case int(qt.Key_Escape):
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
	p.popup.Resize(p.conf.ScanPopupWidth, p.conf.ScanPopupHeight)
}
