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

type ArticleView struct {
	*widgets.QTextBrowser

	doQuery func(string)
}

func NewArticleView() *ArticleView {
	widget := widgets.NewQTextBrowser(nil)
	// widget := webengine.NewQWebEngineView(nil)
	widget.SetReadOnly(true)
	widget.SetOpenExternalLinks(true)
	widget.SetOpenLinks(false)
	return &ArticleView{
		QTextBrowser: widget,
	}
}

func (view *ArticleView) SetupCustomHandlers() {
	doQuery := view.doQuery
	if doQuery == nil {
		panic("doQuery is not set")
	}
	mediaPlayer := multimedia.NewQMediaPlayer(nil, 0)

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
	rightClickOnWord := ""

	view.ConnectContextMenuEvent(func(event *gui.QContextMenuEvent) {
		event.Ignore()
		// menu := webview.CreateStandardContextMenu2(event.GlobalPos())
		menu := view.CreateStandardContextMenu()
		// actions := menu.Actions()
		// log.Println("actions", actions)
		// menu.Actions() panic
		// https://github.com/therecipe/qt/issues/1286
		// firstAction := menu.ActiveAction()

		action := widgets.NewQAction2("Query", view)
		action.ConnectTriggered(func(checked bool) {
			text := view.TextCursor().SelectedText()
			if text != "" {
				doQuery(strings.Trim(text, queryForceTrimChars))
				return
			}
			if rightClickOnWord != "" {
				doQuery(rightClickOnWord)
			}
		})
		menu.InsertAction(nil, action)
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
				rightClickOnWord = strings.Trim(cursor.SelectedText(), punctuation)
				if rightClickOnWord != "" {
					log.Printf("Right-clicked on word %#v\n", rightClickOnWord)
				}
			}
		}
		view.MouseReleaseEventDefault(event)
	})
}
