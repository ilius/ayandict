package application

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/ilius/ayandict/v3/pkg/application/frequency"
	"github.com/ilius/ayandict/v3/pkg/dictmgr"
	common "github.com/ilius/go-dict-commons"
	qt "github.com/mappu/miqt/qt6"
)

const resultFlags = uint32(
	common.ResultFlag_FixAudio |
		common.ResultFlag_FixFileSrc |
		common.ResultFlag_FixWordLink |
		common.ResultFlag_ColorMapping)

type QueryArgs struct {
	ArticleView    *ArticleView
	ResultsLabel   *qt.QLabel
	ResultList     *ResultListWidget
	HeaderLabel    *HeaderLabel
	HistoryView    *HistoryView
	PostQuery      func(string)
	Entry          *qt.QLineEdit
	ModeCombo      *qt.QComboBox
	FrequencyTable *frequency.FrequencyTable
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

func (w *QueryArgs) SetNoResult(query string) {
	w.ArticleView.SetHtml(fmt.Sprintf("No results for %#v", query))
	w.HeaderLabel.SetText("")
	w.AddHistoryAndFrequency(query)
}

func (queryArgs *QueryArgs) onQuery(query string, isAuto bool) {
	if query == "" {
		if !isAuto {
			queryArgs.ArticleView.SetHtml("")
			queryArgs.HeaderLabel.SetText("")
		}
		return
	}
	t := time.Now()
	mode := dictmgr.SearchModeFuzzy
	switch queryArgs.ModeCombo.CurrentIndex() {
	case 1: // StartWith
		mode = dictmgr.SearchModeStartWith
	case 2: // Regex
		if isAuto && !conf.SearchOnTypeOnRegex {
			return
		}
		mode = dictmgr.SearchModeRegex
	case 3: // Glob
		mode = dictmgr.SearchModeGlob
	case 4: // WordMatch
		mode = dictmgr.SearchModeWordMatch
	}
	results := dictmgr.LookupHTML(query, conf, mode, resultFlags, 0)
	slog.Debug("LookupHTML running time", "dt", time.Since(t), "query", query)
	queryArgs.ResultList.SetResults(results)
	queryArgs.ResultsLabel.SetText(fmt.Sprintf("Results: %d", len(results)))
	if len(results) == 0 {
		if !isAuto {
			queryArgs.SetNoResult(query)
		}
	}
	if isAuto {
		if len(results) > 0 {
			if results[0].Score() == 200 {
				queryArgs.AddHistoryAndFrequency(query)
			}
		}
	} else {
		queryArgs.AddHistoryAndFrequency(query)
	}
	queryArgs.PostQuery(query)
}
