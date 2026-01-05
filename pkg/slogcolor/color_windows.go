//go:build windows

package slogcolor

func init() {
	// Opt-in for ansi color support for current process.
	// https://learn.microsoft.com/en-us/windows/console/console-virtual-terminal-sequences#output-sequences
	enabled := true
	EnableColorsStdout(&enabled)
}
