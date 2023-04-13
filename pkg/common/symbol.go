package common

import (
	"fmt"
	"unicode/utf8"
)

func DefaultSymbol(dictName string) string {
	symbol, _ := utf8.DecodeRune([]byte(dictName))
	return fmt.Sprintf("[%s]", string(symbol))
}
