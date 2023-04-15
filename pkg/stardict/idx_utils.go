package stardict

type WordPrefixMap map[rune]map[int]bool

func (wpm WordPrefixMap) Add(term []rune, termIndex int) {
	for _, word := range splitRunes(term, ' ') {
		prefix := word[0]
		m, ok := wpm[prefix]
		if !ok {
			m = map[int]bool{}
			wpm[prefix] = m
		}
		m[termIndex] = true
	}
}
