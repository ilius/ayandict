package dictmgr

import (
	"log/slog"
	"sort"
	"strings"

	"github.com/ilius/ayandict/v3/pkg/config"
	"github.com/ilius/ayandict/v3/pkg/dictmgr/internal/dicts"
	common "github.com/ilius/go-dict-commons"
)

type QueryMode uint8

const (
	QueryModeFuzzy QueryMode = iota
	QueryModeStartWith
	QueryModeRegex
	QueryModeGlob
	QueryModeWordMatch
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
	ds := dicts.DictSettingsMap[dictName]
	if ds == nil {
		ds = &dicts.DictionarySettings{}
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
			slog.Error("error in SearchRegex: " + err.Error())
			return nil
		}
		return results
	case QueryModeGlob:
		if !ds.Glob() {
			return nil
		}
		results, err := dic.SearchGlob(query, workerCount, timeout)
		if err != nil {
			slog.Error("error in SearchGlob: " + err.Error())
			return nil
		}
		return results
	case QueryModeWordMatch:
		if !ds.WordMatch() {
			return nil
		}
		return dic.SearchWordMatch(query, workerCount, timeout)
	}
	if !ds.Fuzzy() {
		return nil
	}
	return dic.SearchFuzzy(query, workerCount, timeout)
}

func sortResults(results []common.SearchResultIface) {
	sort.Slice(results, func(i, j int) bool {
		res1 := results[i]
		res2 := results[j]
		score1 := res1.Score()
		score2 := res2.Score()
		if score1 != score2 {
			return score1 > score2
		}
		term1 := strings.ToLower(res1.Terms()[0])
		term2 := strings.ToLower(res2.Terms()[0])
		if term1 != term2 {
			return term1 < term2
		}
		do1 := dicts.DictsOrder[res1.DictName()]
		do2 := dicts.DictsOrder[res2.DictName()]
		if do1 != do2 {
			return do1 < do2
		}
		// if we do not use entryIndex, the resulting order can be random
		// for entries with same headwords
		// and no need to compare headwords for StarDict when we have entryIndex
		// since they are already sorted in idx file.
		// if we added other formats, maybe we can add a config for this
		return res1.EntryIndex() < res2.EntryIndex()
	})
}

func LookupHTML(
	query string,
	conf *config.Config,
	mode QueryMode,
	resultFlags uint32,
	limit int,
) []common.SearchResultIface {
	results := []common.SearchResultIface{}
	for _, dic := range dicts.DictList {
		if dic.Disabled() || !dic.Loaded() {
			continue
		}
		for _, res := range search(dic, conf, mode, query) {
			results = append(results, NewSearchResult(res, dic, conf, resultFlags))
		}
	}
	if limit == 0 {
		limit = conf.MaxResultsTotal
	}
	if len(results) <= limit {
		sortResults(results)
		return results
	}
	// too many results
	// first only sort by score
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score() > results[j].Score()
	})
	// get rid of extra results
	results = results[:limit]
	// get rid of the bottom score (mainly for consistency)
	minCount := limit * 2 / 3
	lastScore := results[len(results)-1].Score()
	for i := len(results) - 2; i > minCount && results[i].Score() == lastScore; i-- {
		results = results[:i]
	}
	// do the final full sort
	sortResults(results)
	return results
}
