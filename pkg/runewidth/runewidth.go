package runewidth

import (
	"os"

	"github.com/ilius/ayandict/v2/pkg/runewidth/uniseg"
)

//go:generate go run script/generate.go

var (
	// EastAsianWidth will be set true if the current locale is CJK
	EastAsianWidth bool

	// StrictEmojiNeutral should be set false if handle broken fonts
	StrictEmojiNeutral bool = true

	// DefaultCondition is a condition in current locale
	DefaultCondition = &Condition{
		EastAsianWidth:     false,
		StrictEmojiNeutral: true,
	}
)

func init() {
	handleEnv()
}

func handleEnv() {
	env := os.Getenv("RUNEWIDTH_EASTASIAN")
	if env == "" {
		EastAsianWidth = IsEastAsian()
	} else {
		EastAsianWidth = env == "1"
	}
	// update DefaultCondition
	if DefaultCondition.EastAsianWidth != EastAsianWidth {
		DefaultCondition.EastAsianWidth = EastAsianWidth
		if len(DefaultCondition.combinedLut) > 0 {
			DefaultCondition.combinedLut = DefaultCondition.combinedLut[:0]
			CreateLUT()
		}
	}
}

type interval struct {
	first rune
	last  rune
}

type table []interval

func inTables(r rune, ts ...table) bool {
	for _, t := range ts {
		if inTable(r, t) {
			return true
		}
	}
	return false
}

func inTable(r rune, t table) bool {
	if r < t[0].first {
		return false
	}

	bot := 0
	top := len(t) - 1
	for top >= bot {
		mid := (bot + top) >> 1

		switch {
		case t[mid].last < r:
			bot = mid + 1
		case t[mid].first > r:
			top = mid - 1
		default:
			return true
		}
	}

	return false
}

var private = table{
	{0x00E000, 0x00F8FF}, {0x0F0000, 0x0FFFFD}, {0x100000, 0x10FFFD},
}

var nonprint = table{
	{0x0000, 0x001F},
	{0x007F, 0x009F},
	{0x00AD, 0x00AD},
	{0x070F, 0x070F},
	{0x180B, 0x180E},
	{0x200B, 0x200F},
	{0x2028, 0x202E},
	{0x206A, 0x206F},
	{0xD800, 0xDFFF},
	{0xFEFF, 0xFEFF},
	{0xFFF9, 0xFFFB},
	{0xFFFE, 0xFFFF},
}

// Condition have flag EastAsianWidth whether the current locale is CJK or not.
type Condition struct {
	combinedLut        []byte
	EastAsianWidth     bool
	StrictEmojiNeutral bool
}

// RuneWidth returns the number of cells in r.
// See http://www.unicode.org/reports/tr11/
func (c *Condition) RuneWidth(r rune) int {
	if r < 0 || r > 0x10FFFF {
		return 0
	}
	if len(c.combinedLut) > 0 {
		return int(c.combinedLut[r>>1]>>(uint(r&1)*4)) & 3
	}
	// optimized version, verified by TestRuneWidthChecksums()
	if !c.EastAsianWidth {
		switch {
		case r < 0x20:
			return 0
		case (r >= 0x7F && r <= 0x9F) || r == 0xAD: // nonprint
			return 0
		case r < 0x300:
			return 1
		case inTable(r, narrow):
			return 1
		case inTables(r, nonprint, combining):
			return 0
		case inTable(r, doublewidth):
			return 2
		default:
			return 1
		}
	} else {
		switch {
		case inTables(r, nonprint, combining):
			return 0
		case inTable(r, narrow):
			return 1
		case inTables(r, ambiguous, doublewidth):
			return 2
		case !c.StrictEmojiNeutral && inTables(r, ambiguous, emoji, narrow):
			return 2
		default:
			return 1
		}
	}
}

// CreateLUT will create an in-memory lookup table of 557056 bytes for faster operation.
// This should not be called concurrently with other operations on c.
// If options in c is changed, CreateLUT should be called again.
func (c *Condition) CreateLUT() {
	const max = 0x110000
	lut := c.combinedLut
	if len(c.combinedLut) != 0 {
		// Remove so we don't use it.
		c.combinedLut = nil
	} else {
		lut = make([]byte, max/2)
	}
	for i := range lut {
		i32 := int32(i * 2)
		x0 := c.RuneWidth(i32)
		x1 := c.RuneWidth(i32 + 1)
		lut[i] = uint8(x0) | uint8(x1)<<4
	}
	c.combinedLut = lut
}

// StringWidth return width as you can see
func (c *Condition) StringWidth(s string) (width int) {
	g := uniseg.NewGraphemes(s)
	for g.Next() {
		var chWidth int
		for _, r := range g.Runes() {
			chWidth = c.RuneWidth(r)
			if chWidth > 0 {
				break // Our best guess at this point is to use the width of the first non-zero-width rune.
			}
		}
		width += chWidth
	}
	return width
}

// StringWidth return width as you can see
func StringWidth(s string) (width int) {
	return DefaultCondition.StringWidth(s)
}

// CreateLUT will create an in-memory lookup table of 557055 bytes for faster operation.
// This should not be called concurrently with other operations.
func CreateLUT() {
	if len(DefaultCondition.combinedLut) > 0 {
		return
	}
	DefaultCondition.CreateLUT()
}
