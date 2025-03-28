package color

import (
	"github.com/ilius/ayandict/v3/pkg/go-colorable"
)

func init() {
	// Opt-in for ansi color support for current process.
	// https://learn.microsoft.com/en-us/windows/console/console-virtual-terminal-sequences#output-sequences
	enabled := true
	colorable.EnableColorsStdout(&enabled)
}
