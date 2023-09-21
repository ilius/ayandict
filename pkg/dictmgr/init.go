package dictmgr

import (
	"github.com/ilius/ayandict/v2/pkg/config"
	"github.com/ilius/ayandict/v2/pkg/dictmgr/internal/dicts"
)

func InitDicts(conf *config.Config) {
	dicts.InitDicts(conf)
}
