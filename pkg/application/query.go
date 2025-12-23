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

func (w *QueryArgs) AddHistoryAndFrequency(query string) {
	if !conf.HistoryDisable {
		w.HistoryView.Add(query)
	}
	if !conf.MostFrequentDisable {
		w.FrequencyTable.Add(query, 1)
		if conf.MostFrequentAutoSave {
			w.FrequencyTable.SaveNoError()
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

func (app *Application) onQuery(query string, isAuto bool) {
	queryArgs := app.queryArgs
	if query == "" {
		if !isAuto {
			queryArgs.ArticleView.SetHtml("")
			queryArgs.HeaderLabel.SetText("")
		}
		return
	}
	t := time.Now()
	modeDesc := queryArgs.ModeCombo.CurrentText()
	mode, ok := searchModeByDesc[modeDesc]
	if !ok {
		slog.Error("invalid serarch mode", "modeDesc", modeDesc)
	}
	if isAuto {
		switch mode {
		case dictmgr.SearchModeRegex:
			if !conf.SearchOnTypeOnRegex {
				return
			}
		case dictmgr.SearchModeSoundex:
			return
		}
	}
	results := dictmgr.LookupHTML(query, conf, mode, resultFlags, 0)
	slog.Debug("LookupHTML running time", "dt", time.Since(t), "query", query)
	app.SetResults(results)
	if len(results) == 0 {
		if !isAuto {
			app.SetNoResult(query)
		}
	}
	if queryArgs.historyOnQuery(isAuto, results) {
		queryArgs.AddHistoryAndFrequency(query)
	}
	app.postQuery(query)
}

func (q *QueryArgs) historyOnQuery(isAuto bool, results []common.SearchResultIface) bool {
	if q.DisableHistory {
		return false
	}
	if !isAuto {
		return true
	}
	// isAuto=true (search-on-type)
	if len(results) > 0 {
		if results[0].Score() == 200 {
			return true
		}
	}
	return false
}
