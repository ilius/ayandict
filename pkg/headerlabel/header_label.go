package headerlabel

import (
	"html/template"
	"log/slog"
	"strings"

	common "codeberg.org/ilius/go-dict-commons"
	"github.com/ilius/ayandict/v3/pkg/config"
	"github.com/ilius/ayandict/v3/pkg/headerlib"
	"github.com/ilius/ayandict/v3/pkg/utils"
	qt "github.com/mappu/miqt/qt6"
)

const maxMenuWidth = 400

type HeaderLabel struct {
	*qt.QLabel

	result common.SearchResultIface

	text string

	doQuery func(string)

	headerTpl *template.Template
}

func NewHeaderLabel(
	conf *config.Config,
	doQuery func(string),
	headerTpl *template.Template,
) *HeaderLabel {
	qLabel := qt.NewQLabel2()
	qLabel.SetTextInteractionFlags(qt.TextSelectableByMouse)
	// | qt.TextSelectableByKeyboard
	qLabel.SetContentsMargins(0, 0, 0, 0)
	qLabel.SetTextFormat(qt.RichText)
	qLabel.SetWordWrap(conf.HeaderWordWrap)
	qLabel.SetSizePolicy2(qt.QSizePolicy__Expanding, qt.QSizePolicy__Minimum)
	label := &HeaderLabel{
		QLabel:    qLabel,
		headerTpl: headerTpl,
	}
	qLabel.OnContextMenuEvent(func(super func(*qt.QContextMenuEvent), event *qt.QContextMenuEvent) {
		event.Ignore()
		menu := label.createContextMenu(qLabel.SelectedText() != "")
		menu.Popup(event.GlobalPos())
	})
	label.doQuery = doQuery
	return label
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
	header, err := headerlib.GetHeader(label.headerTpl, res, 200)
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
	menuWidth := 0
	fm := menu.FontMetrics()
	updateMenuWidth := func(text string) {
		width := fm.HorizontalAdvance(text)
		if width > menuWidth {
			menuWidth = width
		}
	}
	addActionWithText := func(text string, onTrigger func()) {
		menu.AddActionWithText(text).OnTriggered(onTrigger)
		updateMenuWidth(text)
	}
	if selection {
		addActionWithText("Query Selected", func() {
			text := label.SelectedText()
			if text != "" {
				label.doQuery(strings.Trim(text, utils.QueryForceTrimChars))
			}
		})
		addActionWithText("Copy Selected", func() {
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
		updateMenuWidth("Query: " + term)
	}

	addActionWithText("Copy All (Plaintext)", func() {
		qt.QGuiApplication_Clipboard().SetText2(
			plaintextFromHTML(label.Text()),
			qt.QClipboard__Clipboard,
		)
	})
	addActionWithText("Copy All (HTML)", func() {
		qt.QGuiApplication_Clipboard().SetText2(
			label.Text(),
			qt.QClipboard__Clipboard,
		)
	})

	menuWidth += fm.HorizontalAdvance("M")/2 + 60
	if menuWidth > maxMenuWidth {
		menuWidth = maxMenuWidth
	}
	menu.SetMinimumWidth(menuWidth)

	return menu
}

func plaintextFromHTML(htext string) string {
	doc := qt.NewQTextDocument()
	doc.SetHtml(htext)
	return doc.ToPlainText()
}
