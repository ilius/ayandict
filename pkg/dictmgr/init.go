package dictmgr

import (
	"github.com/ilius/ayandict/v3/pkg/config"
	"github.com/ilius/ayandict/v3/pkg/dictmgr/internal/dicts"
)

func InitDicts(conf *config.Config) {
	dicts.InitDicts(conf)
}
