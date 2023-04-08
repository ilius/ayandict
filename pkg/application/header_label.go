package application

import (
	"bytes"
	"html"
	"log"
	"strings"

	"github.com/ilius/ayandict/pkg/common"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

type HeaderLabel struct {
	*widgets.QLabel

	app *widgets.QApplication

	result common.QueryResult

	doQuery func(string)
}

func CreateHeaderLabel(app *widgets.QApplication) *HeaderLabel {
	qlabel := widgets.NewQLabel(nil, 0)
	qlabel.SetTextInteractionFlags(core.Qt__TextSelectableByMouse)
	// | core.Qt__TextSelectableByKeyboard
	// headerLabel.SetAlignment(core.Qt__AlignVCenter)
	qlabel.SetContentsMargins(20, 0, 0, 0)
	qlabel.SetTextFormat(core.Qt__RichText)
	qlabel.SetWordWrap(true)
	qlabel.SetSizePolicy2(expanding, widgets.QSizePolicy__Minimum)
	label := &HeaderLabel{
		QLabel: qlabel,
		app:    app,
	}
	qlabel.ConnectContextMenuEvent(func(event *gui.QContextMenuEvent) {
		event.Ignore()
		menu := label.createContextMenu(qlabel.SelectedText() != "")
		menu.Popup(event.GlobalPos(), nil)
	})
	return label
}

func (label *HeaderLabel) SetResult(res common.QueryResult) {
	label.result = res
	terms := res.Terms()
	termsJoined := html.EscapeString(strings.Join(terms, " | "))
	headerBuf := bytes.NewBuffer(nil)
	err := headerTpl.Execute(headerBuf, HeaderTemplateInput{
		Terms:    terms,
		Term:     termsJoined,
		DictName: res.DictName(),
		Score:    res.Score() / 2,
	})
	if err != nil {
		log.Println(err)
		return
	}
	label.SetText(headerBuf.String())
}

func (label *HeaderLabel) addQueryAction(menu *widgets.QMenu, term string) {
	menu.AddAction("Query: " + term).ConnectTriggered(func(checked bool) {
		res := label.result
		if res == nil {
			return
		}
		label.doQuery(term)
	})
}

func (label *HeaderLabel) createContextMenu(selection bool) *widgets.QMenu {
	menu := widgets.NewQMenu(label.QLabel)
	if selection {
		menu.AddAction("Query Selected").ConnectTriggered(func(checked bool) {
			text := label.SelectedText()
			if text != "" {
				label.doQuery(strings.Trim(text, queryForceTrimChars))
			}
		})
		menu.AddAction("Copy Selected").ConnectTriggered(func(checked bool) {
			text := label.SelectedText()
			if text == "" {
				return
			}
			label.app.Clipboard().SetText(strings.TrimSpace(text), gui.QClipboard__Clipboard)
		})
	}
	terms := label.result.Terms()
	if len(terms) > 10 {
		terms = terms[:10]
	}
	for _, term := range terms {
		label.addQueryAction(menu, term)
	}

	menu.AddAction("Copy All (Plaintext)").ConnectTriggered(func(checked bool) {
		label.app.Clipboard().SetText(
			plaintextFromHTML(label.Text()),
			gui.QClipboard__Clipboard,
		)
	})
	menu.AddAction("Copy All (HTML)").ConnectTriggered(func(checked bool) {
		label.app.Clipboard().SetText(
			label.Text(),
			gui.QClipboard__Clipboard,
		)
	})

	return menu
}
