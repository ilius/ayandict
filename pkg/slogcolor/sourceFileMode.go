package slogcolor

type SourceFileMode int

const (
	// Nop does nothing.
	Nop SourceFileMode = iota

	// ShortFile produces only the filename (for example main.go:69).
	ShortFile

	// LongFile produces the full file path (for example
	// /home/frajer/go/src/myapp/main.go:69).
	LongFile
)
