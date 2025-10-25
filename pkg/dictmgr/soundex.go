package dictmgr

import (
	"log/slog"

	"github.com/ilius/ayandict/v3/pkg/config"
	"github.com/ilius/ayandict/v3/pkg/dictmgr/internal/dicts"
	"github.com/ilius/ayandict/v3/pkg/mysoundex"
	common "github.com/ilius/go-dict-commons"
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
				if delta > 20 {
					delta = 20
				}
				resLow.F_Score = 200 - 2*uint8(delta)
				results = append(results, NewSearchResult(resLow, dic, conf, resultFlags))
			}
		}
	}
	slog.Info("LookupSoundexHTML: got results", "query", query, "count", len(results))
	sortResults(results)
	if limit == 0 {
		limit = conf.MaxResultsTotal
	}
	if len(results) > limit {
		results = results[:limit]
	}
	return results
}
