package wordwrap

import (
	"strings"

	"github.com/ilius/ayandict/v2/pkg/runewidth"
)

func WordWrapByWords(
	words []string,
	limit int,
	trimLeft string,
	trimRight string,
) [][]string {
	if limit <= 0 || len(words) < 2 {
		return [][]string{words}
	}
	firstWord := strings.TrimLeft(words[0], trimLeft)
	lines := [][]string{{firstWord}}
	currentWidth := runewidth.StringWidth(firstWord)
	for _, word := range words[1:] {
		wordWidth := runewidth.StringWidth(word)
		if currentWidth+wordWidth <= limit {
			i := len(lines) - 1
			lines[i] = append(lines[i], word)
			currentWidth += wordWidth
			continue
		}
		if trimRight != "" {
			lineI := len(lines) - 1
			line := lines[lineI]
			wordI := len(line) - 1
			if wordI >= 0 {
				lastWord := line[wordI]
				lastWord = strings.TrimRight(lastWord, trimRight)
				if lastWord == "" {
					lines[lineI] = line[:wordI]
				} else {
					line[wordI] = lastWord
				}
			}
		}
		// going to next line
		t_word := strings.TrimLeft(word, trimLeft)
		if t_word == "" {
			continue
		}
		lines = append(lines, []string{t_word})
		if t_word != word {
			wordWidth = runewidth.StringWidth(t_word)
		}
		currentWidth = wordWidth
	}
	return lines
}
