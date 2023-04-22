package application

import (
	"log"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ilius/ayandict/pkg/mp3duration"
	commons "github.com/ilius/go-dict-commons"
	"github.com/ilius/qt/core"
	"github.com/ilius/qt/gui"
	"github.com/ilius/qt/multimedia"
	"github.com/ilius/qt/widgets"
)

// Source code of QTextEdit says:
// Zooming into HTML documents only works if the font-size is not set to a fixed size.
// https://github.com/qt/qtbase/blob/dev/src/widgets/widgets/qtextedit.cpp#L451
// https://bugreports.qt.io/browse/QTBUG-52751

const (
	startFrag = "<!--StartFragment-->"
	endFrag   = "<!--EndFragment-->"
)

var dummyParagRE = regexp.MustCompile(`<p [^<>]*><br />(</p>|$)`)

type ArticleView struct {
	*widgets.QTextBrowser

	app     *Application
	dpi     float64
	doQuery func(string)

	mediaPlayer *multimedia.QMediaPlayer

	rightClickOnWord string
	rightClickOnUrl  string

	autoPlayMutex sync.Mutex
}

func NewArticleView(app *Application) *ArticleView {
	widget := widgets.NewQTextBrowser(nil)
	// widget := webengine.NewQWebEngineView(nil)
	widget.SetReadOnly(true)
	widget.SetOpenExternalLinks(true)
	widget.SetOpenLinks(false)
	dpi := app.PrimaryScreen().PhysicalDotsPerInch()
	return &ArticleView{
		QTextBrowser: widget,
		app:          app,
		dpi:          dpi,
	}
}

var audioUrlRE = regexp.MustCompile(`href="[^<>"]+\.mp3"`)

func (view *ArticleView) playAudio(qUrl *core.QUrl) {
	log.Println("Playing audio", qUrl.ToString(core.QUrl__PreferLocalFile))
	player := view.mediaPlayer
	content := multimedia.NewQMediaContent2(qUrl)
	player.SetMedia(content, nil)
	player.Play()
}

func (view *ArticleView) autoPlay(text string, count int) {
	if !view.autoPlayMutex.TryLock() {
		return
	}
	defer view.autoPlayMutex.Unlock()
	matches := audioUrlRE.FindAllString(text, count)
	lastIndex := len(matches) - 1
	for index, match := range matches {
		urlStr, err := strconv.Unquote(match[5:])
		if err != nil {
			log.Println(err)
			continue
		}
		qUrl := core.NewQUrl3(urlStr, core.QUrl__TolerantMode)
		// log.Println("Playing audio", urlStr)
		isRemote := qUrl.Scheme() != "file"
		if isRemote {
			qUrlLocal, err := audioCache.Get(urlStr)
			if err != nil {
				log.Println(err)
			} else {
				qUrl = qUrlLocal
				isRemote = false
			}
		}
		view.playAudio(qUrl)
		// log.Println("Duration:", player.Duration())
		// player.Duration() is always zero
		if isRemote {
			time.Sleep(2000 * time.Millisecond)
			continue
		}
		fpath := filePathFromQUrl(qUrl)
		if fpath == "" {
			continue
		}
		// log.Println("Calculating duration of", fpath)
		duration, err := mp3duration.Calculate(fpath)
		if err != nil {
			log.Printf("error in mp3duration.Calculate(%#v): %v", fpath, err)
			continue
		}
		if index < lastIndex {
			duration += conf.AudioAutoPlayWaitBetween
		}
		// log.Println("Sleeping", duration)
		time.Sleep(duration)
	}
}

func (view *ArticleView) SetResult(res commons.SearchResultIface) {
	text := strings.Join(
		res.DefinitionsHTML(),
		"\n<br/>\n",
	)
	text2 := text
	if definitionStyleString != "" {
		text2 = definitionStyleString + text2
	}
	view.SetHtml(text2)
	if conf.Audio && conf.AudioAutoPlay > 0 {
		go view.autoPlay(text, conf.AudioAutoPlay)
	}
}

func (view *ArticleView) createContextMenu() *widgets.QMenu {
	menu := widgets.NewQMenu(view.QTextBrowser)
	menu.AddAction("Query").ConnectTriggered(func(checked bool) {
		text := view.TextCursor().SelectedText()
		if text != "" {
			view.doQuery(strings.Trim(text, queryForceTrimChars))
			return
		}
		if view.rightClickOnWord != "" {
			view.doQuery(view.rightClickOnWord)
		}
	})
	if view.rightClickOnUrl != "" {
		menu.AddAction("Copy Link Target").ConnectTriggered(func(checked bool) {
			view.app.Clipboard().SetText(view.rightClickOnUrl, gui.QClipboard__Clipboard)
		})
	}
	menu.AddAction("Copy").ConnectTriggered(func(checked bool) {
		text := view.TextCursor().SelectedText()
		if text == "" {
			return
		}
		text = strings.TrimSpace(text)
		view.app.Clipboard().SetText(text, gui.QClipboard__Clipboard)
	})
	menu.AddAction("Copy HTML").ConnectTriggered(func(checked bool) {
		text := view.selectedHTML()
		view.app.Clipboard().SetText(text, gui.QClipboard__Clipboard)
	})
	menu.AddAction("Copy All (HTML)").ConnectTriggered(func(checked bool) {
		view.app.Clipboard().SetText(
			view.ToHtml(),
			gui.QClipboard__Clipboard,
		)
	})
	menu.AddAction("Copy All (Plaintext)").ConnectTriggered(func(checked bool) {
		view.app.Clipboard().SetText(
			view.ToPlainText(),
			gui.QClipboard__Clipboard,
		)
	})

	return menu
}

func (view *ArticleView) selectedHTML() string {
	text := view.TextCursor().Selection().ToHtml(core.NewQByteArray2("utf-8", 5))
	body := strings.Index(text, "<body>")
	if body >= 0 {
		text = text[body+6:]
	}
	endBody := strings.Index(text, "</body>")
	if endBody >= 0 {
		text = text[:endBody]
	}
	start := strings.Index(text, startFrag)
	if start >= 0 {
		text = text[start+len(startFrag):]
	}
	end := strings.Index(text, endFrag)
	if end >= 0 {
		text = text[:end]
	}
	text = strings.TrimSpace(text)
	text = dummyParagRE.ReplaceAllString(text, "")
	return text
}

func (view *ArticleView) zoom(delta int) {
	doc := view.Document()
	font := doc.DefaultFont()
	points := fontPointSize(font, view.dpi)
	if points <= 0 {
		log.Printf("bad font size: points=%v, pixels=%v", font.PointSizeF(), font.PixelSize())
		return
	}
	if delta > 0 {
		font.SetPointSizeF(points * conf.ArticleZoomFactor)
	} else {
		font.SetPointSizeF(points / conf.ArticleZoomFactor)
	}
	doc.SetDefaultFont(font)
}

func (view *ArticleView) ZoomIn(ran int) {
	view.zoom(1)
}

func (view *ArticleView) ZoomOut(ran int) {
	view.zoom(-1)
}

func (view *ArticleView) findLinkOnCursor(cursor *gui.QTextCursor) string {
	text := cursor.Selection().ToHtml(core.NewQByteArray2("utf-8", 5))
	// log.Println("findLinkOnCursor:", text)
	start := strings.Index(text, startFrag)
	if start >= 0 {
		text = text[start:]
	}
	start = strings.Index(text, "href=")
	if start < 0 {
		// log.Println("findLinkOnCursor: did not find end href=")
		return ""
	}
	text = text[start+5:]
	end := strings.Index(text[1:], text[:1])
	if end < 0 {
		// log.Println("findLinkOnCursor: did not find end quote")
		return ""
	}
	urlStr, err := strconv.Unquote(text[:end+2])
	if err != nil {
		// log.Println("error in Unquote", err)
		return ""
	}
	return urlStr
}

func (view *ArticleView) SetupCustomHandlers() {
	doQuery := view.doQuery
	if doQuery == nil {
		panic("doQuery is not set")
	}
	mediaPlayer := multimedia.NewQMediaPlayer(nil, 0)
	view.mediaPlayer = mediaPlayer

	copyAction := widgets.NewQAction2("Copy", view)
	view.AddAction(copyAction)
	copyAction.SetShortcut(gui.NewQKeySequence2("Ctrl+C", gui.QKeySequence__PortableText))
	copyAction.ConnectTriggered(func(checked bool) {
		text := view.TextCursor().SelectedText()
		if text == "" {
			return
		}
		text = strings.TrimSpace(text)
		view.app.Clipboard().SetText(text, gui.QClipboard__Clipboard)
	})

	view.ConnectAnchorClicked(func(qUrl *core.QUrl) {
		host := qUrl.Host(core.QUrl__FullyDecoded)
		// log.Printf(
		// 	"AnchorClicked: %#v, host=%#v = %#v\n",
		// 	link.ToString(core.QUrl__None),
		// 	host,
		// 	link.Host(core.QUrl__FullyEncoded),
		// )
		if qUrl.Scheme() == "bword" {
			if host != "" {
				doQuery(host)
			} else {
				log.Printf("AnchorClicked: %#v\n", qUrl.ToString(core.QUrl__None))
			}
			return
		}
		path := qUrl.Path(core.QUrl__FullyDecoded)
		// log.Printf("scheme=%#v, host=%#v, path=%#v", link.Scheme(), host, path)
		switch qUrl.Scheme() {
		case "":
			doQuery(path)
			return
		case "file":
			// log.Printf("host=%#v, ext=%#v", host, ext)
			switch filepath.Ext(path) {
			case ".mp3", ".wav", ".ogg":
				view.playAudio(qUrl)
				return
			}
		case "http", "https":
			// log.Printf("host=%#v, ext=%#v", host, ext)
			switch filepath.Ext(path) {
			case ".mp3", ".wav", ".ogg":
				qUrlLocal, err := audioCache.Get(qUrl.ToString(core.QUrl__None))
				if err != nil {
					log.Println(err)
				} else {
					qUrl = qUrlLocal
				}
				view.playAudio(qUrl)
				return
			}
		}
		gui.QDesktopServices_OpenUrl(qUrl)
	})
	// menuStyleOpt := widgets.NewQStyleOptionMenuItem()
	// style := app.Style()
	// queryMenuIcon := style.StandardIcon(widgets.QStyle__SP_ArrowUp, menuStyleOpt, nil)

	// we set this on right-button MouseRelease when no text is selected
	// and read it when Query is selected from context menu
	// may not be pretty or concurrent-safe! but seems to work!

	view.ConnectContextMenuEvent(func(event *gui.QContextMenuEvent) {
		event.Ignore()
		menu := view.createContextMenu()
		menu.Popup(event.GlobalPos(), nil)
	})
	view.ConnectMouseReleaseEvent(func(event *gui.QMouseEvent) {
		text := view.TextCursor().SelectedText()
		switch event.Button() {
		case core.Qt__MiddleButton:
			if text != "" {
				doQuery(strings.Trim(text, queryForceTrimChars))
			}
			return
		case core.Qt__RightButton:
			if text == "" {
				cursor := view.CursorForPosition(event.Pos())
				cursor.Select(gui.QTextCursor__WordUnderCursor)
				// it doesn't actually select the word in GUI

				urlStr := view.findLinkOnCursor(cursor)
				if urlStr != "" {
					// log.Println("right-click on url:", urlStr)
					view.rightClickOnUrl = urlStr
				}

				view.rightClickOnWord = strings.Trim(cursor.SelectedText(), punctuation)
				if view.rightClickOnWord != "" {
					log.Printf("Right-clicked on word %#v\n", view.rightClickOnWord)
				}
			}
		}
		view.MouseReleaseEventDefault(event)
	})

	view.ConnectWheelEvent(func(event *gui.QWheelEvent) {
		if event.Modifiers()&core.Qt__ControlModifier == 0 {
			view.WheelEventDefault(event)
			return
		}
		delta := event.AngleDelta().Y()
		if delta == 0 {
			return
		}
		view.zoom(delta)
	})
}

func (view *ArticleView) KeyPressEventDefault(event gui.QKeyEvent_ITF) {
	switch event.QKeyEvent_PTR().Key() {
	case int(core.Qt__Key_Up):
		if conf.ArticleArrowKeys {
			view.VerticalScrollBar().TriggerAction(widgets.QAbstractSlider__SliderSingleStepSub)
			return
		}
	case int(core.Qt__Key_Down):
		if conf.ArticleArrowKeys {
			view.VerticalScrollBar().TriggerAction(widgets.QAbstractSlider__SliderSingleStepAdd)
			return
		}
	}
	view.QTextBrowser.KeyPressEventDefault(event)
}
