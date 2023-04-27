package application

import (
	"bytes"
	"html"
	"strings"

	"github.com/ilius/ayandict/v2/pkg/dictmgr"
	"github.com/ilius/ayandict/v2/pkg/qerr"
	common "github.com/ilius/go-dict-commons"
	"github.com/ilius/qt/core"
	"github.com/ilius/qt/gui"
	"github.com/ilius/qt/widgets"
)

type HeaderLabel struct {
	*widgets.QLabel

	app *Application

	result common.SearchResultIface

	text string

	doQuery func(string)
}

func CreateHeaderLabel(app *Application) *HeaderLabel {
	qLabel := widgets.NewQLabel(nil, 0)
	qLabel.SetTextInteractionFlags(core.Qt__TextSelectableByMouse)
	// | core.Qt__TextSelectableByKeyboard
	qLabel.SetContentsMargins(0, 0, 0, 0)
	qLabel.SetTextFormat(core.Qt__RichText)
	qLabel.SetWordWrap(conf.HeaderWordWrap)
	qLabel.SetSizePolicy2(expanding, widgets.QSizePolicy__Minimum)
	label := &HeaderLabel{
		QLabel: qLabel,
		app:    app,
	}
	qLabel.ConnectContextMenuEvent(func(event *gui.QContextMenuEvent) {
		event.Ignore()
		menu := label.createContextMenu(qLabel.SelectedText() != "")
		menu.Popup(event.GlobalPos(), nil)
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
	terms := res.Terms()
	termsJoined := html.EscapeString(strings.Join(terms, " | "))
	headerBuf := bytes.NewBuffer(nil)
	dictName := res.DictName()
	err := headerTpl.Execute(headerBuf, HeaderTemplateInput{
		Terms:     terms,
		Term:      termsJoined,
		DictName:  dictName,
		Score:     res.Score() >> 1,
		ShowTerms: dictmgr.DictShowTerms(dictName),
	})
	if err != nil {
		qerr.Errorf("Error formatting header label: %v", err)
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
