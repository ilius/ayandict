package slogcolor

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
)

var (
	// Output defines the standard output of the print functions. By default,
	// os.Stdout is used.
	Output = NewColorableStdout()

	// Error defines a color supporting writer for os.Stderr.
	Error = NewColorableStderr()

	// colorsCache is used to reduce the count of created Color objects and
	// allows to reuse already created objects with required Attribute.
	colorsCache   = make(map[Attribute]*Color)
	colorsCacheMu sync.Mutex // protects colorsCache
)

// Color defines a custom color object which is defined by SGR parameters.
type Color struct {
	params []Attribute
}

// Attribute defines a single SGR Code
type Attribute int

const escape = "\x1b"

// Base attributes
const (
	Reset Attribute = iota
	Bold
	Faint
)

const (
	ResetBold Attribute = iota + 22
)

// Foreground text colors
const (
	FgBlack Attribute = iota + 30
	FgRed
	FgGreen
	FgYellow
	FgBlue
	FgMagenta
	FgCyan
	FgWhite
)

// Foreground Hi-Intensity text colors
const (
	FgHiBlack Attribute = iota + 90
	FgHiRed
	FgHiGreen
	FgHiYellow
	FgHiBlue
	FgHiMagenta
	FgHiCyan
	FgHiWhite
)

// Background text colors
const (
	BgBlack Attribute = iota + 40
	BgRed
	BgGreen
	BgYellow
	BgBlue
	BgMagenta
	BgCyan
)

// NewColor returns a newly created color object.
func NewColor(value ...Attribute) *Color {
	c := &Color{
		params: make([]Attribute, 0),
	}

	c.Add(value...)
	return c
}

// Add is used to chain SGR parameters. Use as many as parameters to combine
// and create custom color objects. Example: Add(FgRed, Underline).
func (c *Color) Add(value ...Attribute) *Color {
	c.params = append(c.params, value...)
	return c
}

// Sprint is just like Print, but returns a string instead of printing it.
func (c *Color) Sprint(a ...interface{}) string {
	return c.wrap(fmt.Sprint(a...))
}

// Sprintf is just like Printf, but returns a string instead of printing it.
func (c *Color) Sprintf(format string, a ...interface{}) string {
	return c.wrap(fmt.Sprintf(format, a...))
}

// SprintFunc returns a new function that returns colorized strings for the
// given arguments with fmt.Sprint(). Useful to put into or mix into other
// string. Windows users should use this in conjunction with Output, example:
//
//	put := New(FgYellow).SprintFunc()
//	fmt.Fprintf(Output, "This is a %s", put("warning"))
func (c *Color) SprintFunc() func(a ...interface{}) string {
	return func(a ...interface{}) string {
		return c.wrap(fmt.Sprint(a...))
	}
}

// sequence returns a formatted SGR sequence to be plugged into a "\x1b[...m"
// an example output might be: "1;36" -> bold cyan
func (c *Color) sequence() string {
	format := make([]string, len(c.params))
	for i, v := range c.params {
		format[i] = strconv.Itoa(int(v))
	}

	return strings.Join(format, ";")
}

// wrap wraps the s string with the colors attributes. The string is ready to
// be printed.
func (c *Color) wrap(s string) string {
	if c.isNoColorSet() {
		return s
	}

	return c.format() + s + c.unformat()
}

func (c *Color) format() string {
	return fmt.Sprintf("%s[%sm", escape, c.sequence())
}

func (c *Color) unformat() string {
	// return fmt.Sprintf("%s[%dm", escape, Reset)
	// for each element in sequence let's use the specific reset escape, or the generic one if not found
	format := make([]string, len(c.params))
	for i := range c.params {
		format[i] = strconv.Itoa(int(Reset))
	}

	return fmt.Sprintf("%s[%sm", escape, strings.Join(format, ";"))
}

func (c *Color) isNoColorSet() bool {
	return false
}

func getCachedColor(p Attribute) *Color {
	colorsCacheMu.Lock()
	defer colorsCacheMu.Unlock()

	c, ok := colorsCache[p]
	if !ok {
		c = NewColor(p)
		colorsCache[p] = c
	}

	return c
}

func colorString(format string, p Attribute) string {
	c := getCachedColor(p)
	return c.SprintFunc()(format)
}

// HiWhiteString is a convenient helper function to return a string with hi-intensity white
// foreground.
func HiWhiteString(format string) string {
	return colorString(format, FgHiWhite)
}
