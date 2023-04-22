package dictmgr

import (
	"fmt"
	"log"

	"github.com/ilius/ayandict/pkg/qerr"
	commons "github.com/ilius/go-dict-commons"
)

func Hash(info commons.Dictionary) string {
	log.Println("Calculating hash for", info.DictName())
	b_hash, err := info.CalcHash()
	if err != nil {
		qerr.Error(err)
		return ""
	}
	return fmt.Sprintf("%x", b_hash)
}
