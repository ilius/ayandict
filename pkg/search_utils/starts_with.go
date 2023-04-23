package search_utils

import "strings"

func ScoreStartsWith(terms []string, query string) uint8 {
	bestScore := uint8(0)
	for termI, termOrig := range terms {
		term := strings.ToLower(termOrig)
		if !strings.HasPrefix(term, query) {
			continue
		}
		subtract := uint8(3)
		if termI < 3 {
			subtract = uint8(termI)
		}
		deltaLen := len(term) - len(query)
		subtract2 := uint8(20)
		if deltaLen < 20 {
			subtract2 = uint8(deltaLen)
		}
		score := 200 - subtract - subtract2
		if score > bestScore {
			bestScore = score
		}
	}
	return bestScore
}
