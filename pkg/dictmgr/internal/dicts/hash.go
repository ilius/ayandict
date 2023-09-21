package dicts

import (
	"fmt"
	"log"

	"github.com/ilius/ayandict/v2/pkg/qerr"
	common "github.com/ilius/go-dict-commons"
)

func Hash(info common.Dictionary) string {
	log.Println("Calculating hash for", info.DictName())
	b_hash, err := info.CalcHash()
	if err != nil {
		qerr.Error(err)
		return ""
	}
	return fmt.Sprintf("%x", b_hash)
}
