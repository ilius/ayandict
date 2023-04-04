package stardict

import (
	"strings"
	"testing"
)

func Test_splitRunes(t *testing.T) {
	test := func(str string) {
		words := strings.Split(str, " ")
		runeWords := splitRunes([]rune(str), ' ')
		if len(words) != len(runeWords) {
			t.Fatalf("str=%#v, len(words)=%v, len(runeWords)=%v", str, len(words), len(runeWords))
		}
		for i, word := range words {
			if string(runeWords[i]) != word {
				t.Fatalf("word=%#v, string(runeWord)=%v", word, string(runeWords[i]))
			}
		}
	}
	test("")
	test("test 1")
	test("hello world")
	test("hello")
	test("hello world  abc")
}
