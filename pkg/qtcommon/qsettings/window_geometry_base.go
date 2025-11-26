package qsettings

import (
	"encoding/json"
	"log/slog"
	"os"
	"path"
)

type WindowSettings struct {
	X         int  `json:"x"`
	Y         int  `json:"y"`
	Width     int  `json:"w"`
	Height    int  `json:"h"`
	Maximized bool `json:"max"`
}

func (s *WindowSettings) Save(fname string) {
	fpath := path.Join(stateDir, fname+".json")
	b, err := json.Marshal(s)
	if err != nil {
		slog.Error("error encoding window settings", "err", err, "path", fpath)
		return
	}
	err = os.WriteFile(fpath, b, 0o644)
	if err != nil {
		slog.Error("error saving window settings", "err", err, "path", fpath)
	}
}

func (s *WindowSettings) Load(fname string) {
	fpath := path.Join(stateDir, fname+".json")
	b, err := os.ReadFile(fpath)
	if err != nil {
		if !os.IsNotExist(err) {
			slog.Error("error reading window settings", "err", err, "path", fpath)
		}
		return
	}
	err = json.Unmarshal(b, s)
	if err != nil {
		slog.Error("error decoding window settings", "err", err, "path", fpath, "text", string(b))
	}
}
