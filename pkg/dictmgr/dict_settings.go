package dictmgr

import (
	common "github.com/ilius/go-dict-commons"
)

type DictSettings struct {
	Symbol string `json:"symbol"`
	Order  int    `json:"order"`
	Hash   string `json:"hash"`

	HideTermsHeader bool `json:"terms_header"`
}

func NewDictSettings(dic common.Dictionary, index int) *DictSettings {
	return &DictSettings{
		Symbol: common.DefaultSymbol(dic.DictName()),
		Order:  index,
		Hash:   "",

		HideTermsHeader: false,
	}
}
