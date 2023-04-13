package dictmgr

import (
	"github.com/ilius/ayandict/pkg/common"
	"github.com/ilius/ayandict/pkg/config"
	"github.com/ilius/ayandict/pkg/stardict"
)

func LookupHTML(
	query string,
	conf *config.Config,
) []common.QueryResult {
	return stardict.LookupHTML(
		query,
		conf,
		dictsOrder,
	)
}
