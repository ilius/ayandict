package search_utils

import "strings"

type ScoreFuzzyArgs struct {
	Query          string
	QueryRunes     []rune
	QueryMainWord  []rune
	QueryWordCount int
	MinWordCount   int
	MainWordIndex  int
}

func ScoreFuzzy(
	terms []string,
	args *ScoreFuzzyArgs,
) uint8 {
	bestScore := uint8(0)
	for termI, termOrig := range terms {
		subtract := uint8(3)
		if termI < 3 {
			subtract = uint8(termI)
		}
		term := strings.ToLower(termOrig)
		if term == args.Query {
			return 200 - subtract
		}
		words := strings.Split(term, " ")
		if len(words) < args.MinWordCount {
			continue
		}
		score := Similarity(args.QueryRunes, []rune(term), subtract)
		if score > bestScore {
			bestScore = score
			if score >= 180 {
				continue
			}
		}
		if len(words) > 1 {
			bestWordScore := uint8(0)
			for wordI, word := range words {
				wordScore := Similarity(args.QueryMainWord, []rune(word), subtract)
				if wordScore < 50 {
					continue
				}
				if wordI == args.MainWordIndex {
					wordScore -= 1
				} else {
					wordScore -= wordScore / 10
				}
				if wordScore > bestWordScore {
					bestWordScore = wordScore
				}
			}
			if bestWordScore < 50 {
				continue
			}
			if args.QueryWordCount > 1 {
				bestWordScore = bestWordScore>>1 + bestWordScore/7
			}
			if bestWordScore > bestScore {
				bestScore = bestWordScore
			}
		}
	}
	return bestScore
}
