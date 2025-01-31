package dictmgr

import (
	"log/slog"
	"math/rand"
	"sort"

	"github.com/ilius/ayandict/v3/pkg/config"
	"github.com/ilius/ayandict/v3/pkg/dictmgr/internal/dicts"
	common "github.com/ilius/go-dict-commons"
)

func entryCount(dic common.Dictionary) int {
	if dic.Disabled() {
		return 0
	}
	n, err := dic.EntryCount()
	if err != nil {
		slog.Error("error", "err", err)
	}
	return n
}

func RandomEntry(conf *config.Config, resultFlags uint32) *SearchResult {
	dn := len(dicts.DictList)
	sums := make([]int, dn+1)
	for i, dic := range dicts.DictList {
		sums[i+1] = sums[i] + entryCount(dic)
	}
	totalEntryN := sums[dn]
	if totalEntryN == 0 {
		return nil
	}
	totalEntryI := rand.Intn(totalEntryN)
	dicIndex := sort.Search(dn, func(i int) bool {
		return totalEntryI < sums[i+1]
	})
	dic := dicts.DictList[dicIndex]
	relEntryI := totalEntryI - sums[dicIndex]
	slog.Debug("RandomEntry", "index", relEntryI, "dictName", dic.DictName())
	entry := dic.EntryByIndex(relEntryI)
	if entry == nil {
		slog.Error("ENTRY NOT FOUND", "index", relEntryI, "dictName", dic.DictName())
		return nil
	}
	entry.F_Score = 200
	return NewSearchResult(entry, dic, conf, resultFlags)
}
