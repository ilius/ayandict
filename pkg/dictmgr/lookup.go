package dictmgr

import (
	"sort"

	"github.com/ilius/ayandict/pkg/common"
	"github.com/ilius/ayandict/pkg/config"
)

func LookupHTML(
	query string,
	conf *config.Config,
) []common.SearchResultIface {
	results := []common.SearchResultIface{}
	for _, dic := range dicList {
		if dic.Disabled() || !dic.Loaded() {
			continue
		}
		for _, res := range dic.SearchFuzzy(query) {
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
