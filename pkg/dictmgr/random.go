package dictmgr

import (
	"log"
	"math/rand"
	"sort"

	"github.com/ilius/ayandict/v2/pkg/config"
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

func RandomEntry(conf *config.Config) *SearchResult {
	dn := len(dicList)
	sums := make([]int, dn+1)
	for i, dic := range dicList {
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
	dic := dicList[dicIndex]
	relEntryI := totalEntryI - sums[dicIndex]
	log.Printf("RandomEntry: %v from %v", relEntryI, dic.DictName())
	entry := dic.EntryByIndex(relEntryI)
	if entry == nil {
		log.Printf("ENTRY NOT FOUND: index %v in %v", relEntryI, dic.DictName())
		return nil
	}
	entry.F_Score = 200
	return &SearchResult{
		SearchResultLow: entry,
		dic:             dic,
		conf:            conf,
	}
}
