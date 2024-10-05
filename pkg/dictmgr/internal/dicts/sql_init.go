//go:build !nosql

package dicts

import (
	"github.com/ilius/ayandict/v2/pkg/qtcommon/qerr"
	sqldict "github.com/ilius/go-dict-sql"
)

func init() {
	sqldict.ErrorHandler = func(err error) {
		qerr.Error(err)
	}
	sqldictOpen = sqldict.Open
}
