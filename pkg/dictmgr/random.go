package dictmgr

import (
	"log"
	"math/rand"
	"sort"

	"github.com/ilius/ayandict/v2/pkg/config"
	"github.com/ilius/ayandict/v2/pkg/dictmgr/internal/dicts"
	common "github.com/ilius/go-dict-commons"
)

func entryCount(dic common.Dictionary) int {
	if dic.Disabled() {
		return 0
	}
	n, err := dic.EntryCount()
	if err != nil {
		log.Println(err)
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
	log.Printf("RandomEntry: %v from %v", relEntryI, dic.DictName())
	entry := dic.EntryByIndex(relEntryI)
	if entry == nil {
		log.Printf("ENTRY NOT FOUND: index %v in %v", relEntryI, dic.DictName())
		return nil
	}
	entry.F_Score = 200
	return NewSearchResult(entry, dic, conf, resultFlags)
}
