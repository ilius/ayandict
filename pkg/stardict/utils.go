package stardict

import "unicode"

func splitRunes(s []rune, c rune) [][]rune {
	if len(s) == 0 {
		return [][]rune{nil}
	}
	res := [][]rune{}
	var buf []rune
	lastPos := 0
	for i, x := range s {
		if x == c {
			if len(buf) > 0 {
				res = append(res, buf)
			}
			buf = nil
			lastPos = i + 1
			continue
		}
		buf = s[lastPos : i+1]
	}
	if len(buf) > 0 {
		res = append(res, buf)
	}
	return res
}

func lowerRunes(s []rune) []rune {
	s2 := make([]rune, len(s))
	for i, c := range s {
		s2[i] = unicode.ToLower(c)
	}
	return s2
}

func stringListFromRunes(ls [][]rune) []string {
	sls := make([]string, len(ls))
	for i, s := range ls {
		sls[i] = string(s)
	}
	return sls
}
