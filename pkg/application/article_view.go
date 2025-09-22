package application

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ilius/ayandict/v3/pkg/dictmgr"
	"github.com/ilius/ayandict/v3/pkg/mp3duration"
	common "github.com/ilius/go-dict-commons"
	qt "github.com/mappu/miqt/qt6"
	"github.com/mappu/miqt/qt6/multimedia" // not in latest tag
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
	*qt.QTextBrowser

	app     *Application
	dpi     float64
	doQuery func(string)

	mediaPlayer *multimedia.QMediaPlayer

	autoPlayMutex sync.Mutex

	dictName string
}

func NewArticleView(app *Application) *ArticleView {
	widget := qt.NewQTextBrowser(nil)
	// widget := webengine.NewQWebEngineView(nil)
	widget.SetReadOnly(true)
	widget.SetOpenExternalLinks(true)
	widget.SetOpenLinks(false)
	dpi := qt.QGuiApplication_PrimaryScreen().PhysicalDotsPerInch()
	view := &ArticleView{
		QTextBrowser: widget,
		app:          app,
		dpi:          dpi,
	}
	widget.OnKeyPressEvent(func(super func(*qt.QKeyEvent), event *qt.QKeyEvent) {
		view.onKeyPressEvent(event)
		super(event)
	})
	return view
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

func (view *ArticleView) playAudio(qUrl *qt.QUrl) {
	urlStr := qUrl.ToLocalFile()
	slog.Info("Playing audio", "url", urlStr)
	if conf.AudioMPV && view.playAudioMPV(urlStr) {
		return
	}
	player := view.mediaPlayer
	player.SetSource(qUrl)
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
		qUrl := qt.NewQUrl4(urlStr, qt.QUrl__TolerantMode)
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

func (view *ArticleView) onContextMenuEvent(_ func(event *qt.QContextMenuEvent), event *qt.QContextMenuEvent) {
	event.Ignore()
	menu := view.createContextMenu(event.Pos())
	menu.Popup(event.GlobalPos())
}

func (view *ArticleView) createContextMenuWithSelection(menu *qt.QMenu, selected string) {
	trimmed := strings.Trim(selected, queryForceTrimChars)
	if trimmed != "" {
		selected = trimmed
	}
	label := "Query Selection"
	if len(selected) < 20 {
		label = "Query: " + selected
	}
	menu.AddActionWithText(label).OnTriggered(func() {
		view.doQuery(selected)
	})
	menu.AddActionWithText("Copy Selection").OnTriggered(func() {
		qt.QGuiApplication_Clipboard().SetText2(strings.TrimSpace(selected), qt.QClipboard__Clipboard)
	})
	menu.AddActionWithText("Copy Selection HTML").OnTriggered(func() {
		qt.QGuiApplication_Clipboard().SetText2(view.selectedHTML(), qt.QClipboard__Clipboard)
	})
}

func (view *ArticleView) createContextMenuNoSelection(menu *qt.QMenu, pos *qt.QPoint) {
	cursor := view.CursorForPosition(pos)
	cursor.Select(qt.QTextCursor__WordUnderCursor) // it doesn't actually select the word in GUI

	cursorUrl := view.findLinkOnCursor(cursor)
	if cursorUrl != "" {
		menu.AddActionWithText("Copy Link Target").OnTriggered(func() {
			qt.QGuiApplication_Clipboard().SetText2(cursorUrl, qt.QClipboard__Clipboard)
		})
	}

	cursorWord := strings.Trim(cursor.SelectedText(), punctuation)
	if cursorWord != "" {
		slog.Debug("Right-clicked on word", "word", fmt.Sprintf("%#v", cursorWord))
		menu.AddActionWithText("Query: " + cursorWord).OnTriggered(func() {
			view.doQuery(cursorWord)
		})
	}
}

func (view *ArticleView) createContextMenu(pos *qt.QPoint) *qt.QMenu {
	menu := qt.NewQMenu(view.QTextBrowser.QWidget)

	selected := view.TextCursor().SelectedText()

	if selected == "" {
		view.createContextMenuNoSelection(menu, pos)
	} else {
		view.createContextMenuWithSelection(menu, selected)
	}

	menu.AddActionWithText("Copy All (HTML)").OnTriggered(func() {
		qt.QGuiApplication_Clipboard().SetText2(
			view.ToHtml(),
			qt.QClipboard__Clipboard,
		)
	})
	menu.AddActionWithText("Copy All (Plaintext)").OnTriggered(func() {
		qt.QGuiApplication_Clipboard().SetText2(
			view.ToPlainText(),
			qt.QClipboard__Clipboard,
		)
	})

	return menu
}

func (view *ArticleView) selectedHTML() string {
	text := view.TextCursor().Selection().ToHtml()
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

func (view *ArticleView) findLinkOnCursor(cursor *qt.QTextCursor) string {
	text := cursor.Selection().ToHtml()
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
	view.OnAnchorClicked(func(qUrl *qt.QUrl) {
		host := qUrl.Host()
		if qUrl.Scheme() == "bword" {
			if host != "" {
				view.doQuery(host)
			} else {
				slog.Debug("AnchorClicked", "url", qUrl.ToString())
			}
			return
		}
		path := qUrl.Path()
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
				qUrlLocal, err := audioCache.Get(qUrl.ToString()) // qt.QUrl__None
				if err != nil {
					slog.Error("error", "err", err)
				} else {
					qUrl = qUrlLocal
				}
				view.playAudio(qUrl)
				return
			}
		}
		qt.QDesktopServices_OpenUrl(qUrl)
	})
}

func (view *ArticleView) setupMouseReleaseEvent() {
	view.OnMouseReleaseEvent(func(super func(*qt.QMouseEvent), event *qt.QMouseEvent) {
		switch event.Button() {
		case qt.MiddleButton:
			selected := view.TextCursor().SelectedText()
			if selected != "" {
				view.doQuery(strings.Trim(selected, queryForceTrimChars))
			}
			return
		}
		super(event)
	})
}

func (view *ArticleView) setupWheelEvent() {
	view.OnWheelEvent(func(super func(*qt.QWheelEvent), event *qt.QWheelEvent) {
		if event.Modifiers()&qt.ControlModifier == 0 {
			super(event)
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

	mediaPlayer := multimedia.NewQMediaPlayer()
	view.mediaPlayer = mediaPlayer

	copyAction := qt.NewQAction5("Copy", view.QObject)
	view.AddAction(copyAction)
	copyAction.SetShortcut(qt.NewQKeySequence7("Ctrl+C", qt.QKeySequence__PortableText))
	copyAction.OnTriggered(func() {
		text := view.TextCursor().SelectedText()
		if text == "" {
			return
		}
		text = strings.TrimSpace(text)
		qt.QGuiApplication_Clipboard().SetText2(text, qt.QClipboard__Clipboard)
	})

	view.setupAnchorClicked()

	// menuStyleOpt := qt.NewQStyleOptionMenuItem()
	// style := app.Style()
	// queryMenuIcon := style.StandardIcon(qt.QStyle__SP_ArrowUp, menuStyleOpt, nil)

	// we set this on right-button MouseRelease when no text is selected
	// and read it when Query is selected from context menu
	// may not be pretty or concurrent-safe! but seems to work!

	view.OnContextMenuEvent(view.onContextMenuEvent)
	view.setupMouseReleaseEvent()
	view.setupWheelEvent()
}

func (view *ArticleView) onKeyPressEvent(event *qt.QKeyEvent) {
	switch event.Key() {
	case int(qt.Key_Up):
		if conf.ArticleArrowKeys {
			view.VerticalScrollBar().TriggerAction(qt.QAbstractSlider__SliderSingleStepSub)
			return
		}
	case int(qt.Key_Down):
		if conf.ArticleArrowKeys {
			view.VerticalScrollBar().TriggerAction(qt.QAbstractSlider__SliderSingleStepAdd)
			return
		}
	}
}
