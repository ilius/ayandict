package stardict

import (
	"strings"
	"unicode/utf8"
)

type WordPrefixMap map[rune]map[int]bool

func (wpm WordPrefixMap) Add(term string, termIndex int) {
	for _, word := range strings.Split(strings.ToLower(term), " ") {
		prefix, _ := utf8.DecodeRuneInString(word)
		m, ok := wpm[prefix]
		if !ok {
			m = map[int]bool{}
			wpm[prefix] = m
		}
		m[termIndex] = true
	}
}
