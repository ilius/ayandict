package dictmgr

import (
	"log/slog"

	"github.com/ilius/ayandict/v3/pkg/config"
	"github.com/ilius/ayandict/v3/pkg/dictmgr/internal/dicts"
	"github.com/ilius/ayandict/v3/pkg/mysoundex"
	common "github.com/ilius/go-dict-commons"
	"github.com/ilius/go-dict-commons/search_utils"
)

var soundexSearcher *mysoundex.SoundexSearcher

func SetSoundexSearcher(ss *mysoundex.SoundexSearcher) {
	soundexSearcher = ss
}

func LookupSoundexHTML(
	query string,
	conf *config.Config,
	resultFlags uint32,
	limit int,
) []common.SearchResultIface {
	if soundexSearcher == nil {
		slog.Error("soundex is not enabled")
		return nil
	}
	results := []common.SearchResultIface{}
	workerCount := conf.SearchWorkerCount
	timeout := conf.SearchTimeout
	args := &search_utils.ScoreFuzzyArgs{
		Query:          query,
		QueryRunes:     []rune(query),
		QueryWordCount: 1,
	}
	for _, word := range soundexSearcher.Lookup(query) {
		for _, dic := range dicts.DictList {
			if dic.Disabled() || !dic.Loaded() {
				continue
			}
			tmpResults := dic.SearchExact(word, workerCount, timeout)
			for _, resLow := range tmpResults {
				if len(resLow.F_Terms) == 0 {
					slog.Warn("bad result with no terms", "resLow", resLow)
					continue
				}
				delta := len(query) - len(resLow.F_Terms[0])
				if delta < 0 {
					delta = -delta
				}
				if delta > 10 {
					delta = 10
				}
				score := search_utils.ScoreFuzzy(
					resLow.F_Terms,
					args,
					nil,
				)
				if score == 0 {
					score = score + 10 - uint8(delta)
				}
				resLow.F_Score = score
				results = append(results, NewSearchResult(resLow, dic, conf, resultFlags))
			}
		}
	}
	slog.Info("LookupSoundexHTML: got results", "query", query, "count", len(results))
	sortResults(results)
	if limit == 0 {
		limit = int(conf.MaxResultsTotal)
	}
	if len(results) > limit {
		results = results[:limit]
	}
	return results
}
