package uniseg

import "unicode/utf8"

// The number of bits to shift the boundary information returned by [Step] to
// obtain the monospace width of the grapheme cluster.
const ShiftWidth = 4

// The bit positions by which states are shifted by the [Step] function. These
// values must ensure state values defined for each of the boundary algorithms
// don't overlap (and that they all still fit in a single int). These must
// correspond to the Mask constants.
const (
	shiftPropState = 21 // No mask as these are always the remaining bits.
)

// The bit mask used to extract the state returned by the [Step] function, after
// shifting. These values must correspond to the shift constants.
const (
	maskGraphemeState = 0xf
)

// Step returns the first grapheme cluster (user-perceived character) found in
// the given byte slice. It also returns information about the boundary between
// that grapheme cluster and the one following it as well as the monospace width
// of the grapheme cluster. There are three types of boundary information: word
// boundaries, sentence boundaries, and line breaks. This function is therefore
// a combination of [FirstGraphemeCluster], [FirstWord], [FirstSentence], and
// [FirstLineSegment].
//
// The "boundaries" return value can be evaluated as follows:
//
//   - boundaries >> ShiftWidth: The width of the grapheme cluster for most
//     monospace fonts where a value of 1 represents one character cell.
//
// This function can be called continuously to extract all grapheme clusters
// from a byte slice, as illustrated in the examples below.
//
// If you don't know which state to pass, for example when calling the function
// for the first time, you must pass -1. For consecutive calls, pass the state
// and rest slice returned by the previous call.
//
// The "rest" slice is the sub-slice of the original byte slice "b" starting
// after the last byte of the identified grapheme cluster. If the length of the
// "rest" slice is 0, the entire byte slice "b" has been processed. The
// "cluster" byte slice is the sub-slice of the input slice containing the
// first identified grapheme cluster.
//
// Given an empty byte slice "b", the function returns nil values.
//
// While slightly less convenient than using the Graphemes class, this function
// has much better performance and makes no allocations. It lends itself well to
// large byte slices.
//
// Note that in accordance with [UAX #14 LB3], the final segment will end with
// a mandatory line break (boundaries&MaskLine == LineMustBreak). You can choose
// to ignore this by checking if the length of the "rest" slice is 0 and calling
// [HasTrailingLineBreak] or [HasTrailingLineBreakInString] on the last rune.
//
// [UAX #14 LB3]: https://www.unicode.org/reports/tr14/#Algorithm
// StepString is like [Step] but its input and outputs are strings.
func StepString(str string, state int) (cluster, rest string, boundaries int, newState int) {
	// An empty byte slice returns nothing.
	if len(str) == 0 {
		return
	}

	// Extract the first rune.
	r, length := utf8.DecodeRuneInString(str)
	if len(str) <= length { // If we're already past the end, there is nothing else to parse.
		prop := property(graphemeCodePoints, r)
		return str, "", runeWidth(r, prop) << ShiftWidth, grAny
	}

	// If we don't know the state, determine it now.
	var graphemeState, firstProp int
	remainder := str[length:]
	if state < 0 {
		graphemeState, firstProp, _ = transitionGraphemeState(state, r)
	} else {
		graphemeState = state & maskGraphemeState
		firstProp = state >> shiftPropState
	}

	// Transition until we find a grapheme cluster boundary.
	width := runeWidth(r, firstProp)
	for {
		var (
			graphemeBoundary bool
			prop             int
		)

		r, l := utf8.DecodeRuneInString(remainder)
		remainder = str[length+l:]

		graphemeState, prop, graphemeBoundary = transitionGraphemeState(graphemeState, r)

		if graphemeBoundary {
			boundary := (width << ShiftWidth)
			return str[:length], str[length:], boundary, graphemeState | (prop << shiftPropState)
		}

		if r == vs16 {
			width = 2
		} else if firstProp != prExtendedPictographic && firstProp != prRegionalIndicator && firstProp != prL {
			width += runeWidth(r, prop)
		} else if firstProp == prExtendedPictographic {
			if r == vs15 {
				width = 1
			} else {
				width = 2
			}
		}

		length += l
		if len(str) <= length {
			return str, "", (width << ShiftWidth), grAny | (prop << shiftPropState)
		}
	}
}
