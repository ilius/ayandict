package dictmgr

import (
	"sort"

	"github.com/ilius/ayandict/pkg/common"
	"github.com/ilius/ayandict/pkg/config"
	"github.com/ilius/ayandict/pkg/qerr"
)

type QueryMode uint8

const (
	QueryModeFuzzy QueryMode = iota
	QueryModeStartWith
	QueryModeRegex
	QueryModeGlob
)

func search(dic common.Dictionary, mode QueryMode, query string) []*common.SearchResultLow {
	switch mode {
	case QueryModeStartWith:
		return dic.SearchStartWith(query)
	case QueryModeRegex:
		results, err := dic.SearchRegex(query)
		if err != nil {
			qerr.Error(err)
			return nil
		}
		return results
	case QueryModeGlob:
		results, err := dic.SearchGlob(query)
		if err != nil {
			qerr.Error(err)
			return nil
		}
		return results
	}
	return dic.SearchFuzzy(query)
}

func LookupHTML(
	query string,
	conf *config.Config,
	mode QueryMode,
) []common.SearchResultIface {
	results := []common.SearchResultIface{}

	for _, dic := range dicList {
		if dic.Disabled() || !dic.Loaded() {
			continue
		}
		for _, res := range search(dic, mode, query) {
			results = append(results, &SearchResult{
				SearchResultLow: res,
				dic:             dic,
				conf:            conf,
			})
		}
	}
	sort.Slice(results, func(i, j int) bool {
		res1 := results[i]
		res2 := results[j]
		score1 := res1.Score()
		score2 := res2.Score()
		if score1 != score2 {
			return score1 > score2
		}
		return dictsOrder[res1.DictName()] < dictsOrder[res2.DictName()]
	})
	cutoff := conf.MaxResultsTotal
	if cutoff > 0 && len(results) > cutoff {
		results = results[:cutoff]
	}
	return results
}
