//go:build !nosql

package dicts

import (
	sqldict "codeberg.org/ilius/go-dict-sql"
	"github.com/ilius/ayandict/v2/pkg/qtcommon/qerr"
)

func init() {
	sqldict.ErrorHandler = func(err error) {
		qerr.Error(err)
	}
	sqldictOpen = sqldict.Open
}
