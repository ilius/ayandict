package main

import (
	"os"

	qt "github.com/mappu/miqt/qt6"
)

type SearchableQTextBrowser struct {
	Browser     *qt.QTextBrowser
	Widget      *qt.QWidget
	frame       *qt.QFrame
	searchEntry *qt.QLineEdit
	searchFrame *qt.QFrame
	docCursor   *qt.QTextCursor
}

func (st *SearchableQTextBrowser) init() {
	mainFrame := st.frame
	mainLayout := qt.NewQVBoxLayout2()
	mainFrame.SetLayout(mainLayout.Layout())
	textBrowser := st.Browser
	searchFrame := st.searchFrame
	searchEntry := st.searchEntry

	// default ContentsMargins:
	// QFrame: 0, QVBoxLayout: 11, QTextBrowser: 1
	mainLayout.SetContentsMargins(0, 0, 0, 0)

	searchLayout := qt.NewQHBoxLayout2()
	searchFrame.SetLayout(searchLayout.Layout())

	searchEntry.SetPlaceholderText("Find text...")
	searchEntry.OnReturnPressed(st.findNext)
	searchEntry.OnTextChanged(st.highlightAll)

	findButton := qt.NewQPushButton3("Find Next")
	findButton.OnClicked(st.findNext)

	searchLayout.AddWidget2(searchEntry.QWidget, 1)
	searchLayout.AddWidget2(findButton.QWidget, 0)

	mainLayout.AddWidget3(textBrowser.QWidget, 1, 0)
	mainLayout.AddWidget3(searchFrame.QWidget, 0, 0)

	searchFrame.SetVisible(false)

	// --- Shortcuts ---
	qt.NewQShortcut2(
		qt.NewQKeySequence2("Ctrl+F"),
		textBrowser.QObject,
	).OnActivated(st.showBar)

	qt.NewQShortcut2(
		qt.NewQKeySequence2("Esc"),
		textBrowser.QObject,
	).OnActivated(st.hideBar)
}

func (st *SearchableQTextBrowser) clearHighlights() {
	doc := st.Browser.Document()
	cursor := qt.NewQTextCursor2(doc)
	cursor.BeginEditBlock()
	format := qt.NewQTextCharFormat()
	cursor.Select(qt.QTextCursor__Document)
	cursor.SetCharFormat(format)
	cursor.EndEditBlock()
}

func (st *SearchableQTextBrowser) highlightAll(query string) {
	st.clearHighlights()
	if query == "" {
		return
	}

	doc := st.Browser.Document()
	cursor := qt.NewQTextCursor2(doc)
	palette := st.Browser.Palette()

	// Theme-aware colors
	allBg := palette.ColorWithCr(qt.QPalette__Highlight)
	allBg.SetAlphaF(0.8)

	currentFg := palette.ColorWithCr(qt.QPalette__HighlightedText)

	formatAll := qt.NewQTextCharFormat()
	formatAll.SetBackground(qt.NewQBrush3(allBg))

	formatCurrent := qt.NewQTextCharFormat()
	formatCurrent.SetForeground(qt.NewQBrush3(currentFg))
	formatCurrent.SetFontWeight(int(qt.QFont__Bold))

	// Highlight all matches
	cursor = qt.NewQTextCursor2(doc)
	for {
		cursor = doc.Find2(query, cursor)
		if cursor.IsNull() {
			break
		}
		cursor.MergeCharFormat(formatAll)
	}

	// Highlight current match with bright text
	activeCursor := st.Browser.TextCursor()
	if !activeCursor.IsNull() && activeCursor.HasSelection() {
		activeCursor.MergeCharFormat(formatCurrent)
	}
}

func (st *SearchableQTextBrowser) findNext() {
	query := st.searchEntry.Text()
	if query == "" {
		return
	}

	flags := qt.QTextDocument__FindFlag(0)
	found := st.Browser.Find2(query, flags)
	if !found {
		st.docCursor.MovePosition3(qt.QTextCursor__Start, qt.QTextCursor__MoveAnchor, 0)
		st.Browser.SetTextCursor(st.docCursor)
		st.Browser.Find2(query, flags)
	}

	st.highlightAll(query)
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

	browser.SetPlainText(string(textB))

	window := qt.NewQWidget(nil)
	window.SetWindowTitle("QTextBrowser Search Example")
	window.Resize(600, 400)
	layout := qt.NewQVBoxLayout2()
	window.SetLayout(layout.Layout())
	layout.AddWidget(stb.Widget)

	window.Show()
	_ = qt.QApplication_Exec()
}
