package dictmgr

import (
	"sort"

	"github.com/ilius/ayandict/v2/pkg/config"
	"github.com/ilius/ayandict/v2/pkg/qerr"
	common "github.com/ilius/go-dict-commons"
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
	dictName := dic.DictName()
	ds := dictSettingsMap[dictName]
	if ds == nil {
		ds = &DictSettings{}
	}
	switch mode {
	case QueryModeStartWith:
		if !ds.StartWith() {
			return nil
		}
		return dic.SearchStartWith(query, workerCount, timeout)
	case QueryModeRegex:
		if !ds.Regex() {
			return nil
		}
		results, err := dic.SearchRegex(query, workerCount, timeout)
		if err != nil {
			qerr.Error(err)
			return nil
		}
		return results
	case QueryModeGlob:
		if !ds.Glob() {
			return nil
		}
		results, err := dic.SearchGlob(query, workerCount, timeout)
		if err != nil {
			qerr.Error(err)
			return nil
		}
		return results
	}
	if !ds.Fuzzy() {
		return nil
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
		// if we do not use entryIndex, the resulting order can be random
		// for entries with same headwords
		// and no need to compare headwords for StarDict when we have entryIndex
		// since they are already sorted in idx file.
		// if we added other formats, maybe we can add a config for this
		// term1 := res1.Terms()[0]
		// term2 := res2.Terms()[0]
		// if term1 != term2 {
		// 	return term1 < term2
		// }
		return res1.EntryIndex() < res2.EntryIndex()
	})
	cutoff := conf.MaxResultsTotal
	if cutoff > 0 && len(results) > cutoff {
		results = results[:cutoff]
	}
	return results
}
