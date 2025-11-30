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
	searchEntry.OnTextChanged(view.onSearchEntryChange)

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

func (view *SearchableQTextBrowser) findAll(query string) []*qt.QTextCursor {
	doc := view.Browser.Document()
	cursor := qt.NewQTextCursor2(doc)
	cursors := []*qt.QTextCursor{}
	for {
		cursor = doc.Find2(query, cursor)
		if cursor.IsNull() {
			break
		}
		cursors = append(cursors, cursor)
	}
	return cursors
}

func (view *SearchableQTextBrowser) onSearchEntryChange(query string) {
	view.highlightAll(query)
}

func (view *SearchableQTextBrowser) clearHighlights() {
	for _, cc := range view.lastFormats {
		cc.Cursor.SetCharFormat(cc.Format)
		cc.Cursor.ClearSelection()
	}
	activeCursor := view.Browser.TextCursor()
	if !activeCursor.IsNull() {
		activeCursor.ClearSelection()
	}
	view.lastFormats = nil
}

func (view *SearchableQTextBrowser) highlightAll(query string) {
	view.clearHighlights()
	if query == "" {
		return
	}

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
	activeCursor := view.Browser.TextCursor()

	cursors := view.findAll(query)
	lastFormats := []CursorAndCharFormat{}
	for _, cursor := range cursors {
		lastFormats = append(lastFormats, CursorAndCharFormat{
			Cursor: cursor,
			Format: cursor.CharFormat(),
		})
	}
	if !activeCursor.IsNull() {
		lastFormats = append(lastFormats, CursorAndCharFormat{
			Cursor: activeCursor,
			Format: activeCursor.CharFormat(),
		})
	}
	for _, cursor := range cursors {
		cursor.MergeCharFormat(formatAll)
	}
	if !activeCursor.IsNull() && activeCursor.HasSelection() {
		// Emphasize the current match
		activeCursor.MergeCharFormat(formatCurrent)
	}

	view.lastFormats = lastFormats
}

func (view *SearchableQTextBrowser) findNext() {
	query := view.searchEntry.Text()
	if query == "" {
		return
	}

	flags := qt.QTextDocument__FindFlag(0)
	flags |= qt.QTextDocument__FindCaseSensitively // TODO: add a checkbox
	found := view.Browser.Find2(query, flags)
	if !found {
		cursor := qt.NewQTextCursor2(view.Browser.Document())
		cursor.MovePosition3(qt.QTextCursor__Start, qt.QTextCursor__MoveAnchor, 0)
		view.Browser.SetTextCursor(cursor)
		view.Browser.Find2(query, flags)
	}

	view.highlightAll(query)
	view.Browser.EnsureCursorVisible()
}

func (view *SearchableQTextBrowser) showBar() {
	query := view.searchEntry.Text()

	activeCursor := view.Browser.TextCursor()
	if !activeCursor.IsNull() && activeCursor.HasSelection() {
		query = activeCursor.SelectedText()
		view.searchEntry.SetText(query)
	}

	view.searchFrame.SetVisible(true)
	view.searchEntry.SetFocus()
	view.searchEntry.SelectAll()
	view.highlightAll(query)
}

func (view *SearchableQTextBrowser) hideBar() {
	view.searchFrame.SetVisible(false)
	view.clearHighlights()
	view.Browser.SetFocus()
}

func CreateTextBrowserSearchBar(browser *qt.QTextBrowser) *SearchableQTextBrowser {
	frame := qt.NewQFrame2()
	view := &SearchableQTextBrowser{
		frame:       frame,
		Widget:      frame.QWidget,
		Browser:     browser,
		searchEntry: qt.NewQLineEdit(nil),
		searchFrame: qt.NewQFrame(nil),
	}
	view.init()
	return view
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
