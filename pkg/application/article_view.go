package application

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/multimedia"
	"github.com/therecipe/qt/widgets"
)

// Source code of QTextEdit says:
// Zooming into HTML documents only works if the font-size is not set to a fixed size.
// https://github.com/qt/qtbase/blob/dev/src/widgets/widgets/qtextedit.cpp#L451
// https://bugreports.qt.io/browse/QTBUG-52751

type ArticleView struct {
	*widgets.QTextBrowser

	app *widgets.QApplication

	doQuery func(string)

	rightClickOnWord string
}

func fontPointSize(font *gui.QFont, dpi float64) float64 {
	points := font.PointSizeF()
	if points > 0 {
		return points
	}
	pixels := font.PixelSize()
	return float64(pixels) * 72.0 / dpi
}

func NewArticleView(app *widgets.QApplication) *ArticleView {
	widget := widgets.NewQTextBrowser(nil)
	// widget := webengine.NewQWebEngineView(nil)
	widget.SetReadOnly(true)
	widget.SetOpenExternalLinks(true)
	widget.SetOpenLinks(false)

	dpi := app.PrimaryScreen().PhysicalDotsPerInch()

	widget.ConnectWheelEvent(func(event *gui.QWheelEvent) {
		if event.Modifiers()&core.Qt__ControlModifier == 0 {
			widget.WheelEventDefault(event)
			return
		}
		doc := widget.Document()
		font := doc.DefaultFont()
		delta := event.AngleDelta().Y()
		// log.Println("WheelEvent", font.PixelSize(), font.PointSizeF())
		if delta == 0 {
			return
		}
		points := fontPointSize(font, dpi)
		if points <= 0 {
			log.Printf("bad font size: points=%v, pixels=%v", font.PointSizeF(), font.PixelSize())
			return
		}
		if delta > 0 {
			font.SetPointSizeF(points * conf.WheelZoomFactor)
		} else {
			font.SetPointSizeF(points / conf.WheelZoomFactor)
		}
		doc.SetDefaultFont(font)
	})
	return &ArticleView{
		QTextBrowser: widget,
		app:          app,
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
	menu.AddAction("Copy").ConnectTriggered(func(checked bool) {
		text := view.TextCursor().SelectedText()
		if text == "" {
			return
		}
		text = strings.TrimSpace(text)
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

func (view *ArticleView) SetupCustomHandlers() {
	doQuery := view.doQuery
	if doQuery == nil {
		panic("doQuery is not set")
	}
	mediaPlayer := multimedia.NewQMediaPlayer(nil, 0)

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

	view.ConnectAnchorClicked(func(link *core.QUrl) {
		host := link.Host(core.QUrl__FullyDecoded)
		// log.Printf(
		// 	"AnchorClicked: %#v, host=%#v = %#v\n",
		// 	link.ToString(core.QUrl__None),
		// 	host,
		// 	link.Host(core.QUrl__FullyEncoded),
		// )
		if link.Scheme() == "bword" {
			if host != "" {
				doQuery(host)
			} else {
				log.Printf("AnchorClicked: %#v\n", link.ToString(core.QUrl__None))
			}
			return
		}
		path := link.Path(core.QUrl__FullyDecoded)
		// log.Printf("scheme=%#v, host=%#v, path=%#v", link.Scheme(), host, path)
		switch link.Scheme() {
		case "":
			doQuery(path)
			return
		case "file", "http", "https":
			// log.Printf("host=%#v, ext=%#v", host, ext)
			switch filepath.Ext(path) {
			case ".wav", ".mp3", ".ogg":
				log.Println("Playing audio", link.ToString(core.QUrl__None))
				mediaPlayer.SetMedia(multimedia.NewQMediaContent2(link), nil)
				mediaPlayer.Play()
				return
			}
		}
		gui.QDesktopServices_OpenUrl(link)
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
				view.rightClickOnWord = strings.Trim(cursor.SelectedText(), punctuation)
				if view.rightClickOnWord != "" {
					log.Printf("Right-clicked on word %#v\n", view.rightClickOnWord)
				}
			}
		}
		view.MouseReleaseEventDefault(event)
	})
}
