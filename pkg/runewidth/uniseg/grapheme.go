package uniseg

// Graphemes implements an iterator over Unicode grapheme clusters, or
// user-perceived characters. While iterating, it also provides information
// about word boundaries, sentence boundaries, line breaks, and monospace
// character widths.
//
// After constructing the class via [NewGraphemes] for a given string "str",
// [Graphemes.Next] is called for every grapheme cluster in a loop until it
// returns false. Inside the loop, information about the grapheme cluster as
// well as boundary information and character width is available via the various
// methods (see examples below).
//
// Using this class to iterate over a string is convenient but it is much slower
// than using this package's [Step] or [StepString] functions or any of the
// other specialized functions starting with "First".
type Graphemes struct {
	// The original string.
	original string

	// The remaining string to be parsed.
	remaining string

	// The current grapheme cluster.
	cluster string

	// The byte offset of the current grapheme cluster relative to the original
	// string.
	offset int

	// The current boundary information of the [Step] parser.
	boundaries int

	// The current state of the [Step] parser.
	state int
}

// NewGraphemes returns a new grapheme cluster iterator.
func NewGraphemes(str string) *Graphemes {
	return &Graphemes{
		original:  str,
		remaining: str,
		state:     -1,
	}
}

// Next advances the iterator by one grapheme cluster and returns false if no
// clusters are left. This function must be called before the first cluster is
// accessed.
func (g *Graphemes) Next() bool {
	if len(g.remaining) == 0 {
		// We're already past the end.
		g.state = -2
		g.cluster = ""
		return false
	}
	g.offset += len(g.cluster)
	g.cluster, g.remaining, g.boundaries, g.state = StepString(g.remaining, g.state)
	return true
}

// Runes returns a slice of runes (code points) which corresponds to the current
// grapheme cluster. If the iterator is already past the end or [Graphemes.Next]
// has not yet been called, nil is returned.
func (g *Graphemes) Runes() []rune {
	if g.state < 0 {
		return nil
	}
	return []rune(g.cluster)
}
