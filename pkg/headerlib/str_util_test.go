package headerlib

import (
	"strings"
	"testing"
)

func Test_joinWithMaxLen(t *testing.T) {
	test := func(maxLen int, sep string, out string, strs ...string) {
		actual, usedCount := joinWithMaxLen(strs, sep, maxLen)
		if actual != out {
			t.Fatalf("expected %#v, actual %#v", out, actual)
		}
		if strings.Join(strs[:usedCount], sep) != out {
			t.Fatalf("bad usedStrs with %d items", usedCount)
		}
	}
	test(0, "|", "a", "a")
	test(20, "|", "a", "a")
	test(20, " | ", "a", "a")
	test(20, " | ", "a | b", "a", "b")
	test(1, " | ", "a", "a", "b")
	test(2, " | ", "a", "a", "b")
	test(3, " | ", "a", "a", "b")
	test(4, " | ", "a", "a", "b")
	test(5, " | ", "a | b", "a", "b")
	test(5, " | ", "a | b", "a", "b", "c")
	test(5, " | ", "a | b", "a", "b", "c", "d")
	test(6, " | ", "a | b", "a", "b", "c")
	test(7, " | ", "a | b", "a", "b", "c")
	test(8, " | ", "a | b", "a", "b", "c")
	test(9, " | ", "a | b | c", "a", "b", "c")
}
