package common

import (
	"fmt"
	"log"

	"github.com/ilius/ayandict/pkg/qerr"
)

func Hash(info Dictionary) string {
	log.Println("Calculating hash for", info.DictName())
	b_hash, err := info.CalcHash()
	if err != nil {
		qerr.Error(err)
		return ""
	}
	return fmt.Sprintf("%x", b_hash)
}
