package qsettings

import (
	"encoding/json"
	"log/slog"
	"os"
	"path"
	"time"

	"github.com/ilius/ayandict/v3/pkg/config"
)

const (
	QS_mainwindow = "mainwindow"

	QS_columnwidth = "columnwidth"
)

func saveJson(value any, mainKey string) {
	fpath := path.Join(config.GetConfigDir(), mainKey+".json")
	b, err := json.Marshal(value)
	if err != nil {
		slog.Error("error encoding splitter sizes", "value", value)
		return
	}
	err = os.WriteFile(fpath, b, 0o644)
	if err != nil {
		slog.Error("error writing splitter sizes")
		return
	}
}

func loadJsonIntSlice(mainKey string) []int {
	fpath := path.Join(config.GetConfigDir(), mainKey+".json")
	b, err := os.ReadFile(fpath)
	if err != nil {
		if !os.IsNotExist(err) {
			slog.Error("error reading splitter sizes", "err", err, "path", fpath)
		}
		return nil
	}
	value := []int{}
	err = json.Unmarshal(b, &value)
	if err != nil {
		slog.Error("error decoding splitter sizes", "err", err, "path", fpath, "data", b)
		return nil
	}
	return value
}

func ActionSaveLoop(ch <-chan time.Time, callable func()) {
	var lastSave time.Time
	for {
		lastEvent := <-ch
	Loop1:
		for {
			select {
			case t := <-ch:
				lastEvent = t
			case <-time.After(500 * time.Millisecond):
				break Loop1
			}
		}
		if lastEvent.After(lastSave) {
			callable()
			lastSave = time.Now()
		}
	}
}
