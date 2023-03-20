package stardict

import (
	stardict "github.com/ilius/go-stardict"
	"github.com/ozeidan/fuzzy-patricia/patricia"
)

var trie = patricia.NewTrie()

func BuildFuzzyTrie(dicList []*stardict.Dictionary) {
	for _, dic := range dicList {
		dic.IterKeywords(func(keyword string) {
			trie.Insert(patricia.Prefix(keyword), dic)
		})
	}
}
