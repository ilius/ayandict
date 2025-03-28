package application

import (
	"log/slog"
	"strings"

	"github.com/ilius/ayandict/v3/pkg/headerlib"
	common "github.com/ilius/go-dict-commons"
	qt "github.com/mappu/miqt/qt6"
)

type HeaderLabel struct {
	*qt.QLabel

	app *Application

	result common.SearchResultIface

	text string

	doQuery func(string)
}

func CreateHeaderLabel(app *Application) *HeaderLabel {
	qLabel := qt.NewQLabel2()
	qLabel.SetTextInteractionFlags(qt.TextSelectableByMouse)
	// | qt.TextSelectableByKeyboard
	qLabel.SetContentsMargins(0, 0, 0, 0)
	qLabel.SetTextFormat(qt.RichText)
	qLabel.SetWordWrap(conf.HeaderWordWrap)
	qLabel.SetSizePolicy2(expanding, qt.QSizePolicy__Minimum)
	label := &HeaderLabel{
		QLabel: qLabel,
		app:    app,
	}
	qLabel.OnContextMenuEvent(func(super func(*qt.QContextMenuEvent), event *qt.QContextMenuEvent) {
		event.Ignore()
		menu := label.createContextMenu(qLabel.SelectedText() != "")
		menu.Popup(event.GlobalPos())
	})
	return label
}

func (label *HeaderLabel) ReloadConfig() {
	label.SetWordWrap(conf.HeaderWordWrap)
}

func (label *HeaderLabel) SetText(text string) {
	if text == label.text {
		return
	}
	label.QLabel.SetText(text)
	label.text = text
	// label.QLabel.AdjustSize()
	parent := label.QLabel.ParentWidget()
	parent.AdjustSize()
}

func (label *HeaderLabel) SetResult(res common.SearchResultIface) {
	label.result = res
	header, err := headerlib.GetHeader(headerTpl, res)
	if err != nil {
		slog.Error("error formatting header label: " + err.Error())
		return
	}
	label.SetText(header)
}

func (label *HeaderLabel) addQueryAction(menu *qt.QMenu, term string) {
	menu.AddActionWithText("Query: " + term).OnTriggered(func() {
		res := label.result
		if res == nil {
			return
		}
		label.doQuery(term)
	})
}

func (label *HeaderLabel) createContextMenu(selection bool) *qt.QMenu {
	menu := qt.NewQMenu(label.QLabel.QWidget)
	if selection {
		menu.AddActionWithText("Query Selected").OnTriggered(func() {
			text := label.SelectedText()
			if text != "" {
				label.doQuery(strings.Trim(text, queryForceTrimChars))
			}
		})
		menu.AddActionWithText("Copy Selected").OnTriggered(func() {
			text := label.SelectedText()
			if text == "" {
				return
			}
			qt.QGuiApplication_Clipboard().SetText2(strings.TrimSpace(text), qt.QClipboard__Clipboard)
		})
	}
	terms := label.result.Terms()
	if len(terms) > 10 {
		terms = terms[:10]
	}
	for _, term := range terms {
		label.addQueryAction(menu, term)
	}

	menu.AddActionWithText("Copy All (Plaintext)").OnTriggered(func() {
		qt.QGuiApplication_Clipboard().SetText2(
			plaintextFromHTML(label.Text()),
			qt.QClipboard__Clipboard,
		)
	})
	menu.AddActionWithText("Copy All (HTML)").OnTriggered(func() {
		qt.QGuiApplication_Clipboard().SetText2(
			label.Text(),
			qt.QClipboard__Clipboard,
		)
	})

	return menu
}
