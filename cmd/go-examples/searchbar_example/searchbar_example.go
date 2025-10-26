package main

import (
	"os"

	qt "github.com/mappu/miqt/qt6"
)

type CursorAndCharFormat struct {
	Cursor *qt.QTextCursor
	Format *qt.QTextCharFormat
}

type SearchableQTextBrowser struct {
	Browser     *qt.QTextBrowser
	Widget      *qt.QWidget
	frame       *qt.QFrame
	searchEntry *qt.QLineEdit
	searchFrame *qt.QFrame
	docCursor   *qt.QTextCursor
	lastFormats []CursorAndCharFormat
}

func (view *SearchableQTextBrowser) init() {
	mainFrame := view.frame
	mainLayout := qt.NewQVBoxLayout2()
	mainFrame.SetLayout(mainLayout.Layout())
	textBrowser := view.Browser
	searchFrame := view.searchFrame
	searchEntry := view.searchEntry

	// default ContentsMargins:
	// QFrame: 0, QVBoxLayout: 11, QTextBrowser: 1
	mainLayout.SetContentsMargins(0, 0, 0, 0)

	searchLayout := qt.NewQHBoxLayout2()
	searchFrame.SetLayout(searchLayout.Layout())

	searchEntry.SetPlaceholderText("Find text...")
	searchEntry.OnReturnPressed(view.findNext)
	searchEntry.OnTextChanged(view.highlightAll)

	findButton := qt.NewQPushButton3("Find Next")
	findButton.OnClicked(view.findNext)

	searchLayout.AddWidget2(searchEntry.QWidget, 1)
	searchLayout.AddWidget2(findButton.QWidget, 0)

	mainLayout.AddWidget3(textBrowser.QWidget, 1, 0)
	mainLayout.AddWidget3(searchFrame.QWidget, 0, 0)

	searchFrame.SetVisible(false)

	// --- Shortcuts ---
	qt.NewQShortcut2(
		qt.NewQKeySequence2("Ctrl+F"),
		textBrowser.QObject,
	).OnActivated(view.showBar)

	qt.NewQShortcut2(
		qt.NewQKeySequence2("Esc"),
		textBrowser.QObject,
	).OnActivated(view.hideBar)
}

func (view *SearchableQTextBrowser) clearHighlights() {
	for _, cc := range view.lastFormats {
		cc.Cursor.SetCharFormat(cc.Format)
	}
	view.lastFormats = nil
}

func (view *SearchableQTextBrowser) highlightAll(query string) {
	view.clearHighlights()
	if query == "" {
		return
	}

	doc := view.Browser.Document()
	palette := view.Browser.Palette()

	// Colors from theme
	allBg := palette.ColorWithCr(qt.QPalette__Highlight)
	allBg.SetAlphaF(0.8)

	currentFg := palette.ColorWithCr(qt.QPalette__HighlightedText)

	formatAll := qt.NewQTextCharFormat()
	formatAll.SetBackground(qt.NewQBrush3(allBg))

	formatCurrent := qt.NewQTextCharFormat()
	formatCurrent.SetForeground(qt.NewQBrush3(currentFg))
	formatCurrent.SetFontWeight(int(qt.QFont__Bold))

	// Apply new highlights
	cursor := qt.NewQTextCursor2(doc)

	lastFormats := []CursorAndCharFormat{}
	for {
		cursor = doc.Find2(query, cursor)
		if cursor.IsNull() {
			break
		}
		lastFormats = append(lastFormats, CursorAndCharFormat{
			Cursor: cursor,
			Format: cursor.CharFormat(),
		})
		cursor.MergeCharFormat(formatAll)
	}
	view.lastFormats = lastFormats

	// Emphasize the current match
	activeCursor := view.Browser.TextCursor()
	if !activeCursor.IsNull() && activeCursor.HasSelection() {
		lastFormats = append(lastFormats, CursorAndCharFormat{
			Cursor: activeCursor,
			Format: activeCursor.CharFormat(),
		})
		activeCursor.MergeCharFormat(formatCurrent)
	}
}

func (st *SearchableQTextBrowser) findNext() {
	query := st.searchEntry.Text()
	if query == "" {
		return
	}

	flags := qt.QTextDocument__FindFlag(0)
	flags |= qt.QTextDocument__FindCaseSensitively // TODO: add a checkbox
	found := st.Browser.Find2(query, flags)
	if !found {
		st.docCursor.MovePosition3(qt.QTextCursor__Start, qt.QTextCursor__MoveAnchor, 0)
		st.Browser.SetTextCursor(st.docCursor)
		st.Browser.Find2(query, flags)
	}

	st.highlightAll(query)
	st.Browser.EnsureCursorVisible()
}

func (st *SearchableQTextBrowser) showBar() {
	st.searchFrame.SetVisible(true)
	st.searchEntry.SetFocus()
	st.searchEntry.SelectAll()
	st.highlightAll(st.searchEntry.Text())
}

func (st *SearchableQTextBrowser) hideBar() {
	st.searchFrame.SetVisible(false)
	st.clearHighlights()
	st.Browser.SetFocus()
}

func CreateTextBrowserSearchBar(browser *qt.QTextBrowser) *SearchableQTextBrowser {
	frame := qt.NewQFrame2()
	st := &SearchableQTextBrowser{
		frame:       frame,
		Widget:      frame.QWidget,
		Browser:     browser,
		searchEntry: qt.NewQLineEdit(nil),
		searchFrame: qt.NewQFrame(nil),
		docCursor:   qt.NewQTextCursor2(browser.Document()),
	}
	st.init()
	return st
}

func main() {
	textB, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	_ = qt.NewQApplication(os.Args[1:])

	browser := qt.NewQTextBrowser(nil)
	stb := CreateTextBrowserSearchBar(browser)

	browser.SetHtml(string(textB))

	window := qt.NewQWidget(nil)
	window.SetWindowTitle("QTextBrowser Search Example")
	window.Resize(600, 400)
	layout := qt.NewQVBoxLayout2()
	window.SetLayout(layout.Layout())
	layout.AddWidget(stb.Widget)

	window.Show()
	_ = qt.QApplication_Exec()
}
