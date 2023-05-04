package wordwrap

import (
	"testing"

	"github.com/ilius/is/v2"
)

// TestByWordsNoLimit tests that a limit of zero or less is considered as an infinite
// maximum length.
func TestByWordsNoLimit(t *testing.T) {
	is := is.New(t)
	in := []string{"this", " ", "is", " ", "a", " ", "test", " ", "string"}
	out := WordWrapByWords(in, 0, "", "")
	if !is.Equal(len(out), 1) {
		return
	}
	is.Equal(out[0], in)
}

// TestLimitGreaterThanString tests that a limit that is greater than the input
// string will not break the string up.
func TestByWordsLimitGreaterThanString(t *testing.T) {
	is := is.New(t)
	in := []string{"this", " ", "is", " ", "a", " ", "test", " ", "string"}
	out := WordWrapByWords(in, 100, "", "")
	if !is.Equal(len(out), 1) {
		return
	}
	is.Equal(out[0], in)
}

// TestSpacesWithLimit tests that a word with a limit and spaces will be
// wrapped at spaces to fit in that limit.
func TestByWordsSpaces(t *testing.T) {
	is := is.New(t)
	in := []string{"this", " ", "is", " ", "a", " ", "test", " ", "string"}
	out := WordWrapByWords(in, 6, "", "")
	expected := [][]string{
		{"this", " "},
		{"is", " ", "a", " "},
		{"test", " "},
		{"string"},
	}
	is = is.Msg("expected=%#v, out=%#v", expected, out)
	is.Equal(out, expected)
}

func TestByWordsSpacesTrim(t *testing.T) {
	is := is.New(t)
	in := []string{"this", " ", "is", " ", "a", " ", "test", " ", "string"}
	out := WordWrapByWords(in, 6, " ", " ")
	expected := [][]string{
		{"this"},
		{"is", " ", "a"},
		{"test"},
		{"string"},
	}
	is = is.Msg("expected=%#v, out=%#v", expected, out)
	is.Equal(out, expected)
}

func TestByWordsSpacesTrim2(t *testing.T) {
	is := is.New(t)
	in := []string{"this", " ", "string", " ", "is", " ", "test"}
	out := WordWrapByWords(in, 6, " ", " ")
	expected := [][]string{
		{"this"},
		{"string"},
		{"is"},
		{"test"},
	}
	is = is.Msg("expected=%#v, out=%#v", expected, out)
	is.Equal(out, expected)
}

func TestByWordsSpacesTrim3(t *testing.T) {
	is := is.New(t)
	in := []string{"this", " ", "string", " ", "is", " ", " test"}
	out := WordWrapByWords(in, 6, " ", " ")
	expected := [][]string{
		{"this"},
		{"string"},
		{"is"},
		{"test"},
	}
	is = is.Msg("expected=%#v, out=%#v", expected, out)
	is.Equal(out, expected)
}

func TestByWordsSpacesTrim4(t *testing.T) {
	is := is.New(t)
	in := []string{"this", " ", "string", " ", " test"}
	out := WordWrapByWords(in, 6, " ", " ")
	expected := [][]string{
		{"this"},
		{"string"},
		{"test"},
	}
	is = is.Msg("expected=%#v, out=%#v", expected, out)
	is.Equal(out, expected)
}

func TestByWordsEndOfString(t *testing.T) {
	is := is.New(t)
	in := []string{"12345"}
	out := WordWrapByWords(in, 5, "", "")
	if !is.Equal(len(out), 1) {
		return
	}
	is.Equal(out[0], in)
}
