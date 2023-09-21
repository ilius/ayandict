package dictmgr

import "github.com/ilius/ayandict/v2/pkg/dictmgr/internal/dicts"

func DictSymbol(dictName string) string {
	ds := dicts.DictSettingsMap[dictName]
	if ds == nil {
		return ""
	}
	return ds.Symbol
}

func DictShowTerms(dictName string) bool {
	ds := dicts.DictSettingsMap[dictName]
	if ds == nil {
		return true
	}
	return !ds.HideTermsHeader
}

func CloseDicts() {
	for _, dic := range dicts.DicList {
		if dic.Disabled() {
			continue
		}
		dic.Close()
	}
}
