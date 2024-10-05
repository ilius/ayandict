package dicts

import (
	"fmt"
	"log/slog"

	"github.com/ilius/ayandict/v2/pkg/qtcommon/qerr"
	common "github.com/ilius/go-dict-commons"
)

func Hash(info common.Dictionary) string {
	slog.Info("Calculating hash for", info.DictName())
	b_hash, err := info.CalcHash()
	if err != nil {
		qerr.Error(err)
		return ""
	}
	return fmt.Sprintf("%x", b_hash)
}
