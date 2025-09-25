package dicts

import (
	common "github.com/ilius/go-dict-commons"
)

const (
	FlagNoFuzzy uint16 = 2 << iota
	FlagNoStartWith
	FlagNoRegex
	FlagNoGlob
	FlagNoWordMatch
)

type DictionarySettings struct {
	Symbol string `json:"symbol"`
	Order  int    `json:"order"`
	Hash   string `json:"hash"`

	Flags uint16 `json:"flags,omitempty"`

	HideTermsHeader bool `json:"hide_terms_header,omitempty"`

	AudioVolume int `json:"audio_volume,omitempty"`
}

func (ds *DictionarySettings) Fuzzy() bool {
	return ds.Flags&FlagNoFuzzy == 0
}

func (ds *DictionarySettings) StartWith() bool {
	return ds.Flags&FlagNoStartWith == 0
}

func (ds *DictionarySettings) Regex() bool {
	return ds.Flags&FlagNoRegex == 0
}

func (ds *DictionarySettings) Glob() bool {
	return ds.Flags&FlagNoGlob == 0
}

func (ds *DictionarySettings) WordMatch() bool {
	return ds.Flags&FlagNoWordMatch == 0
}

func (ds *DictionarySettings) SetFuzzy(enable bool) {
	if enable {
		ds.Flags &= ^FlagNoFuzzy
	} else {
		ds.Flags |= FlagNoFuzzy
	}
}

func (ds *DictionarySettings) SetStartWith(enable bool) {
	if enable {
		ds.Flags &= ^FlagNoStartWith
	} else {
		ds.Flags |= FlagNoStartWith
	}
}

func (ds *DictionarySettings) SetRegex(enable bool) {
	if enable {
		ds.Flags &= ^FlagNoRegex
	} else {
		ds.Flags |= FlagNoRegex
	}
}

func (ds *DictionarySettings) SetGlob(enable bool) {
	if enable {
		ds.Flags &= ^FlagNoGlob
	} else {
		ds.Flags |= FlagNoGlob
	}
}

func NewDictSettings(dic common.Dictionary, index int) *DictionarySettings {
	return &DictionarySettings{
		Symbol: common.DefaultSymbol(dic.DictName()),
		Order:  index,
		Hash:   "",

		HideTermsHeader: false,
	}
}
