package logging

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/ilius/ayandict/v3/pkg/config"
	"github.com/ilius/ayandict/v3/pkg/go-color"
	"github.com/ilius/ayandict/v3/pkg/qtcommon/qerr"
	"github.com/ilius/ayandict/v3/pkg/slogcolor"
)

const DefaultLevel = slog.LevelInfo

func SetupGUILogger(noColor bool, level slog.Level) {
	handler := NewColoredHandler(noColor, level)
	slog.SetDefault(slog.New(&CustomHandler{
		Handler: handler,
	}))
}

func NewColoredHandler(noColor bool, level slog.Level) slog.Handler {
	return slogcolor.NewHandler(os.Stdout, &slogcolor.Options{
		Level:         level,
		TimeFormat:    time.DateTime,
		SrcFileMode:   slogcolor.ShortFile,
		SrcFileLength: 0,
		// MsgPrefix:     color.HiWhiteString("| "),
		MsgLength: 0,
		MsgColor:  color.New(),
		NoColor:   noColor,
	})
}

type CustomHandler struct {
	slog.Handler
}

func (h *CustomHandler) showRecordInGUI(record slog.Record) {
	msg := record.Message
	if msg != "" {
		msg = strings.ToUpper(msg[:1]) + msg[1:] // capitalize first character
	}
	// TODO: check how it looks
	attrs := []string{}
	record.Attrs(func(attr slog.Attr) bool {
		attrs = append(attrs, fmt.Sprintf("%s: %v", attr.Key, attr.Value))
		return true
	})
	// \n does not work, <br> and <br/> does
	// <pre> does not work
	msg += "<br>" + strings.Join(attrs, "<br>")
	qerr.Error(msg)
}

func (h *CustomHandler) Handle(ctx context.Context, record slog.Record) error {
	err := h.Handler.Handle(ctx, record)
	if record.Level == slog.LevelError {
		h.showRecordInGUI(record)
	}
	return err
}

func parseLevel(levelStr string) (slog.Level, bool) {
	switch strings.ToLower(levelStr) {
	case "error":
		return slog.LevelError, true
	case "warn", "warning":
		return slog.LevelWarn, true
	case "info":
		return slog.LevelInfo, true
	case "debug":
		return slog.LevelDebug, true
	}
	return slog.LevelInfo, false
}

func SetupLoggerAfterConfigLoad(noColor bool, conf *config.Config) {
	recreateLogger := false
	level := DefaultLevel
	if !noColor && conf.Logging.NoColor {
		noColor = true
		recreateLogger = true
	}
	if conf.Logging.Level != "" {
		configLevel, ok := parseLevel(conf.Logging.Level)
		if ok {
			if configLevel != DefaultLevel {
				level = configLevel
				recreateLogger = true
			}
		} else {
			slog.Error("invalid log level name", "level", conf.Logging.Level)
		}
	}
	if recreateLogger {
		slog.Info("Re-creating logger after loading config")
		SetupGUILogger(noColor, level)
	}
}
