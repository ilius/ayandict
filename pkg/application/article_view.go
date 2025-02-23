package application

import (
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ilius/ayandict/v2/pkg/dictmgr"
	"github.com/ilius/ayandict/v2/pkg/mp3duration"
	common "github.com/ilius/go-dict-commons"
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

	dictName string
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

func (view *ArticleView) playAudioMPV(urlStr string) bool {
	path, err := exec.LookPath("mpv")
	if err != nil {
		slog.Error("error in LookPath", "err", err)
		return false
	}
	args := []string{
		path,
		"--no-video",
		urlStr,
	}
	volume := dictmgr.AudioVolume(view.dictName) * conf.AudioVolume / 100
	args = append(args, "--volume="+strconv.FormatInt(int64(volume), 10))
	cmd := exec.Cmd{
		Path:   path,
		Args:   args,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	err = cmd.Start()
	if err != nil {
		slog.Error("error in mpv: Start", "err", err)
		return false
	}
	return true
}

func (view *ArticleView) playAudio(qUrl *core.QUrl) {
	urlStr := qUrl.ToString(core.QUrl__PreferLocalFile)
	slog.Info("Playing audio", "url", urlStr)
	if conf.AudioMPV && view.playAudioMPV(urlStr) {
		return
	}
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
			slog.Error("error", "err", err)
			continue
		}
		qUrl := core.NewQUrl3(urlStr, core.QUrl__TolerantMode)
		// slog.Info("Playing audio", urlStr)
		isRemote := qUrl.Scheme() != "file"
		if isRemote {
			qUrlLocal, err := audioCache.Get(urlStr)
			if err != nil {
				slog.Error("error", "err", err)
			} else {
				qUrl = qUrlLocal
				isRemote = false
			}
		}
		view.playAudio(qUrl)
		// slog.Info("Duration:", player.Duration())
		// player.Duration() is always zero
		if isRemote {
			time.Sleep(2000 * time.Millisecond)
			continue
		}
		fpath := filePathFromQUrl(qUrl)
		if fpath == "" {
			continue
		}
		// slog.Info("Calculating duration of", fpath)
		duration, err := mp3duration.Calculate(fpath)
		if err != nil {
			slog.Error("error in mp3duration.Calculate", "fpath", fpath, "err", err)
			continue
		}
		if index < lastIndex {
			duration += conf.AudioAutoPlayWaitBetween
		}
		// slog.Info("Sleeping", duration)
		time.Sleep(duration)
	}
}

func (view *ArticleView) SetResult(res common.SearchResultIface) {
	view.dictName = res.DictName()
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

func (view *ArticleView) ZoomIn() {
	doc := view.Document()
	font := doc.DefaultFont()
	points := fontPointSize(font, view.dpi)
	if points <= 0 {
		return
	}
	font.SetPointSizeF(points * conf.ArticleZoomFactor)
	doc.SetDefaultFont(font)
}

func (view *ArticleView) ZoomOut() {
	doc := view.Document()
	font := doc.DefaultFont()
	points := fontPointSize(font, view.dpi)
	if points <= 0 {
		return
	}
	font.SetPointSizeF(points / conf.ArticleZoomFactor)
	doc.SetDefaultFont(font)
}

func (view *ArticleView) findLinkOnCursor(cursor *gui.QTextCursor) string {
	text := cursor.Selection().ToHtml(core.NewQByteArray2("utf-8", 5))
	// slog.Info("findLinkOnCursor:", text)
	start := strings.Index(text, startFrag)
	if start >= 0 {
		text = text[start:]
	}
	start = strings.Index(text, "href=")
	if start < 0 {
		// slog.Info("findLinkOnCursor: did not find end href=")
		return ""
	}
	text = text[start+5:]
	end := strings.Index(text[1:], text[:1])
	if end < 0 {
		// slog.Info("findLinkOnCursor: did not find end quote")
		return ""
	}
	urlStr, err := strconv.Unquote(text[:end+2])
	if err != nil {
		// slog.Error("error in Unquote", err)
		return ""
	}
	return urlStr
}

func (view *ArticleView) setupAnchorClicked() {
	view.ConnectAnchorClicked(func(qUrl *core.QUrl) {
		host := qUrl.Host(core.QUrl__FullyDecoded)
		if qUrl.Scheme() == "bword" {
			if host != "" {
				view.doQuery(host)
			} else {
				slog.Debug("AnchorClicked", "url", qUrl.ToString(core.QUrl__None))
			}
			return
		}
		path := qUrl.Path(core.QUrl__FullyDecoded)
		switch qUrl.Scheme() {
		case "":
			view.doQuery(path)
			return
		case "file":
			switch filepath.Ext(path) {
			case ".mp3", ".wav", ".ogg":
				view.playAudio(qUrl)
				return
			}
		case "http", "https":
			switch filepath.Ext(path) {
			case ".mp3", ".wav", ".ogg":
				qUrlLocal, err := audioCache.Get(qUrl.ToString(core.QUrl__None))
				if err != nil {
					slog.Error("error", "err", err)
				} else {
					qUrl = qUrlLocal
				}
				view.playAudio(qUrl)
				return
			}
		}
		gui.QDesktopServices_OpenUrl(qUrl)
	})
}

func (view *ArticleView) setupMouseReleaseEvent() {
	view.ConnectMouseReleaseEvent(func(event *gui.QMouseEvent) {
		text := view.TextCursor().SelectedText()
		switch event.Button() {
		case core.Qt__MiddleButton:
			if text != "" {
				view.doQuery(strings.Trim(text, queryForceTrimChars))
			}
			return
		case core.Qt__RightButton:
			if text == "" {
				cursor := view.CursorForPosition(event.Pos())
				cursor.Select(gui.QTextCursor__WordUnderCursor)
				// it doesn't actually select the word in GUI

				urlStr := view.findLinkOnCursor(cursor)
				if urlStr != "" {
					// slog.Info("right-click on url:", urlStr)
					view.rightClickOnUrl = urlStr
				}

				view.rightClickOnWord = strings.Trim(cursor.SelectedText(), punctuation)
				if view.rightClickOnWord != "" {
					slog.Debug("Right-clicked on word", "word", view.rightClickOnWord)
				}
			}
		}
		view.MouseReleaseEventDefault(event)
	})
}

func (view *ArticleView) setupWheelEvent() {
	view.ConnectWheelEvent(func(event *gui.QWheelEvent) {
		if event.Modifiers()&core.Qt__ControlModifier == 0 {
			view.WheelEventDefault(event)
			return
		}
		delta := event.AngleDelta().Y()
		if delta == 0 {
			return
		}
		if delta > 0 {
			view.ZoomIn()
		} else {
			view.ZoomOut()
		}
	})
}

// TODO: break down
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

	view.setupAnchorClicked()

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
	view.setupMouseReleaseEvent()
	view.setupWheelEvent()
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
