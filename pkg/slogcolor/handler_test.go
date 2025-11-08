package slogcolor

import (
	"log/slog"
	"os"
	"testing"
	"time"
)

func TestNewHandlerColor(t *testing.T) {
	handler := NewHandler(os.Stdout, &Options{
		Level:       slog.LevelDebug,
		TimeFormat:  time.DateTime,
		SrcFileMode: ShortFile,
		// MsgPrefix:     color.HiWhiteString("| "),
		MsgColor: NewColor(),
	})
	slog.SetDefault(slog.New(handler))
	slog.Debug("test", "a", "b")
	slog.Info("test", "a", "b")
	slog.Warn("test", "a", "b")
	slog.Error("test", "a", "b")
}

func TestNewHandlerNoColor(t *testing.T) {
	handler := NewHandler(os.Stdout, &Options{
		Level:       slog.LevelDebug,
		TimeFormat:  time.DateTime,
		SrcFileMode: ShortFile,
		// MsgPrefix:     color.HiWhiteString("| "),
		MsgColor: NewColor(),
		NoColor:  true,
	})
	slog.SetDefault(slog.New(handler))
	slog.Debug("test", "a", "b")
	slog.Info("test", "a", "b")
	slog.Warn("test", "a", "b")
	slog.Error("test", "a", "b")
}
