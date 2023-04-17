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

func search(
	dic common.Dictionary,
	conf *config.Config,
	mode QueryMode,
	query string,
) []*common.SearchResultLow {
	workerCount := conf.SearchWorkerCount
	timeout := conf.SearchTimeout
	switch mode {
	case QueryModeStartWith:
		return dic.SearchStartWith(query, workerCount, timeout)
	case QueryModeRegex:
		results, err := dic.SearchRegex(query, workerCount, timeout)
		if err != nil {
			qerr.Error(err)
			return nil
		}
		return results
	case QueryModeGlob:
		results, err := dic.SearchGlob(query, workerCount, timeout)
		if err != nil {
			qerr.Error(err)
			return nil
		}
		return results
	}
	return dic.SearchFuzzy(query, workerCount, timeout)
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
		for _, res := range search(dic, conf, mode, query) {
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
		do1 := dictsOrder[res1.DictName()]
		do2 := dictsOrder[res2.DictName()]
		if do1 != do2 {
			return do1 < do2
		}
		return res1.Terms()[0] < res2.Terms()[0]
	})
	cutoff := conf.MaxResultsTotal
	if cutoff > 0 && len(results) > cutoff {
		results = results[:cutoff]
	}
	return results
}
