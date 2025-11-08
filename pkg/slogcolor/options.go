package slogcolor

import (
	"log/slog"
	"time"
)

var DefaultOptions *Options = &Options{
	Level:      slog.LevelInfo,
	TimeFormat: time.DateTime,
	MsgPrefix:  HiWhiteString("| "),
	MsgColor:   NewColor(),
	NoColor:    false,
}

type Options struct {
	// Level reports the minimum level to log.
	// Levels with lower levels are discarded.
	// If nil, the Handler uses [slog.LevelInfo].
	Level slog.Leveler

	// TimeFormat is the time format.
	TimeFormat string

	// MsgPrefix to show prefix before message, default: white colored "| ".
	MsgPrefix string

	// MsgColor is the color of the message, default to empty.
	MsgColor *Color

	// NoColor disables color, default: false.
	NoColor bool
}
