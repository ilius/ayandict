//go:build !nosql

package dictmgr

import (
	"github.com/ilius/ayandict/v2/pkg/qerr"
	sqldict "github.com/ilius/go-dict-sql"
)

func init() {
	sqldict.ErrorHandler = func(err error) {
		qerr.Error(err)
	}
	sqldictOpen = sqldict.Open
}
