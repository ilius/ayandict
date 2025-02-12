//go:build !windows || appengine
// +build !windows appengine

package colorable

import (
	"io"
	"os"
)

// NewColorableStdout returns new instance of Writer which handles escape sequence for stdout.
func NewColorableStdout() io.Writer {
	return os.Stdout
}

// NewColorableStderr returns new instance of Writer which handles escape sequence for stderr.
func NewColorableStderr() io.Writer {
	return os.Stderr
}
