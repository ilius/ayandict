package application

import (
	"fmt"
	"log/slog"
	"time"

	common "codeberg.org/ilius/go-dict-commons"
	"github.com/ilius/ayandict/v3/pkg/articleview"
	"github.com/ilius/ayandict/v3/pkg/dictmgr"
	"github.com/ilius/ayandict/v3/pkg/frequencytable"
	"github.com/ilius/ayandict/v3/pkg/headerlabel"
	qt "github.com/mappu/miqt/qt6"
)

const resultFlags = uint32(
	common.ResultFlag_FixAudio |
		common.ResultFlag_FixFileSrc |
		common.ResultFlag_FixWordLink |
		common.ResultFlag_ColorMapping |
		common.ResultFlag_QTextBrowser)

type QueryArgs struct {
	ArticleView    *articleview.ArticleView
	ResultsLabel   *qt.QLabel
	HeaderLabel    *headerlabel.HeaderLabel
	HistoryView    *HistoryView
	Entry          *qt.QLineEdit
	ModeCombo      *qt.QComboBox
	FrequencyTable *frequencytable.FrequencyTable

	DisableHistory bool // temporarily disable history
}

func (a *QueryArgs) Mode() dictmgr.SearchMode {
	modeDesc := a.ModeCombo.CurrentText()
	mode, ok := searchModeByDesc[modeDesc]
	if !ok {
		slog.Error("invalid serarch mode", "modeDesc", modeDesc)
	}
	return mode
}

func (a *QueryArgs) AddHistoryAndFrequency(query string) {
	if a.DisableHistory {
		return
	}
	if !conf.HistoryDisable {
		a.HistoryView.Add(query)
	}
	if !conf.MostFrequentDisable {
		a.FrequencyTable.Add(query, 1)
		if conf.MostFrequentAutoSave {
			a.FrequencyTable.SaveNoError()
		}
	}
}

func (app *Application) SetNoResult(query string) {
	a := app.queryArgs
	a.ArticleView.SetHtml(fmt.Sprintf("No results for %#v", query))
	a.HeaderLabel.SetText("")
	a.ResultsLabel.SetText("Results: none")
	app.favoriteButton.Hide()
	a.AddHistoryAndFrequency(query)
}

func (app *Application) ResultsLabel() *qt.QLabel {
	return app.queryArgs.ResultsLabel
}

func (app *Application) SetResults(results []common.SearchResultIface) {
	app.resultList.SetResults(results)
}

func (app *Application) onQuery(query string) {
	args := app.queryArgs
	if query == "" {
		args.ArticleView.SetHtml("")
		args.HeaderLabel.SetText("")
		return
	}
	mode := args.Mode()
	startTime := time.Now()
	results := dictmgr.LookupHTML(query, conf, mode, resultFlags, 0)
	slog.Debug("LookupHTML running time", "dt", time.Since(startTime), "query", query)
	app.SetResults(results)
	if len(results) == 0 {
		app.SetNoResult(query)
	}
	args.AddHistoryAndFrequency(query)
	app.postQuery(query)
}

func (app *Application) onQueryAuto(query string) {
	if query == "" {
		return
	}
	mode := app.queryArgs.Mode()
	switch mode {
	case dictmgr.SearchModeRegex:
		if !conf.SearchOnTypeOnRegex {
			return
		}
	case dictmgr.SearchModeSoundex:
		return
	}
	startTime := time.Now()
	results := dictmgr.LookupHTML(query, conf, mode, resultFlags, 0)
	slog.Debug("LookupHTML running time", "dt", time.Since(startTime), "query", query)
	app.SetResults(results)
	if len(results) > 0 && results[0].Score() > 198 {
		app.queryArgs.AddHistoryAndFrequency(query)
	}
	app.postQuery(query)
}
