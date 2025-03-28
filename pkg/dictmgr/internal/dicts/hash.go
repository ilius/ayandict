package dicts

import (
	"fmt"
	"log/slog"

	common "github.com/ilius/go-dict-commons"
)

func Hash(info common.Dictionary) string {
	slog.Info("Calculating dict hash", "dictName", info.DictName())
	b_hash, err := info.CalcHash()
	if err != nil {
		slog.Error("error in CalcHash: " + err.Error())
		return ""
	}
	return fmt.Sprintf("%x", b_hash)
}
