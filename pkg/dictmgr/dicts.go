package dictmgr

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/ilius/ayandict/v2/pkg/dictmgr/internal/dicts"
)

const DictResPathBase = "/dict-res/"

func DictSymbol(dictName string) string {
	ds := dicts.DictSettingsMap[dictName]
	if ds == nil {
		return ""
	}
	return ds.Symbol
}

func DictShowTerms(dictName string) bool {
	ds := dicts.DictSettingsMap[dictName]
	if ds == nil {
		return true
	}
	return !ds.HideTermsHeader
}

func CloseDicts() {
	for _, dic := range dicts.DictList {
		if dic.Disabled() {
			continue
		}
		dic.Close()
	}
}

func DictResFile(dictName string, resPath string) (string, bool) {
	dic, ok := dicts.DictByName[dictName]
	if !ok {
		return "", false
	}
	resDir := dic.ResourceDir()
	if resDir == "" {
		return "", false
	}
	fpath := filepath.Join(resDir, resPath)
	_, err := os.Stat(fpath)
	if err != nil {
		if err != os.ErrNotExist {
			slog.Error("error", "err", err)
		}
		return "", false
	}
	return fpath, true
}

func AudioVolume(dictName string) int {
	ds := dicts.DictSettingsMap[dictName]
	if ds == nil {
		slog.Error("AudioVolume: no Settings value", "dictName", dictName)
		return 200
	}
	if ds.AudioVolume == 0 {
		return 200
	}
	return ds.AudioVolume
}
