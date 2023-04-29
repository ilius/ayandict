package dictmgr

import (
	common "github.com/ilius/go-dict-commons"
)

const (
	FlagNoFuzzy uint16 = 2 << iota
	FlagNoStartWith
	FlagNoRegex
	FlagNoGlob
)

type DictSettings struct {
	Symbol string `json:"symbol"`
	Order  int    `json:"order"`
	Hash   string `json:"hash"`

	Flags uint16 `json:"flags"`

	HideTermsHeader bool `json:"terms_header"`
}

func (ds *DictSettings) Fuzzy() bool {
	return ds.Flags&FlagNoFuzzy == 0
}

func (ds *DictSettings) StartWith() bool {
	return ds.Flags&FlagNoStartWith == 0
}

func (ds *DictSettings) Regex() bool {
	return ds.Flags&FlagNoRegex == 0
}

func (ds *DictSettings) Glob() bool {
	return ds.Flags&FlagNoGlob == 0
}

func (ds *DictSettings) SetFuzzy(enable bool) {
	if enable {
		ds.Flags |= ^FlagNoFuzzy
	} else {
		ds.Flags |= FlagNoFuzzy
	}
}

func (ds *DictSettings) SetStartWith(enable bool) {
	if enable {
		ds.Flags |= ^FlagNoStartWith
	} else {
		ds.Flags |= FlagNoStartWith
	}
}

func (ds *DictSettings) SetRegex(enable bool) {
	if enable {
		ds.Flags |= ^FlagNoRegex
	} else {
		ds.Flags |= FlagNoRegex
	}
}

func (ds *DictSettings) SetGlob(enable bool) {
	if enable {
		ds.Flags |= ^FlagNoGlob
	} else {
		ds.Flags |= FlagNoGlob
	}
}

func NewDictSettings(dic common.Dictionary, index int) *DictSettings {
	return &DictSettings{
		Symbol: common.DefaultSymbol(dic.DictName()),
		Order:  index,
		Hash:   "",

		HideTermsHeader: false,
	}
}
