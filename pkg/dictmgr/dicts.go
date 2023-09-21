package dictmgr

import (
	"log"
	"net/url"
	"os"
	"path/filepath"

	"github.com/ilius/ayandict/v2/pkg/dictmgr/internal/dicts"
	common "github.com/ilius/go-dict-commons"
)

const DictResPathBase = "/dict-res/"

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

func DictResFile(dictName string, resPath string) (string, bool) {
	dic, ok := dicts.DicMap[dictName]
	if !ok {
		return "", false
	}
	resDir := dic.ResourceDir()
	if resDir == "" {
		return "", false
	}
	fpath := filepath.Join(resDir, resPath)
	_, err := os.Stat(fpath)
	if err != nil {
		if err != os.ErrNotExist {
			log.Println(err)
		}
		return "", false
	}
	return fpath, true
}

func dictResURL(dic common.Dictionary, relPath string, flags uint32) string {
	if flags&common.ResultFlag_Web > 0 {
		values := url.Values{}
		values.Add("dictName", dic.DictName())
		values.Add("path", relPath)
		return DictResPathBase + "?" + values.Encode()
	}
	return dic.ResourceURL() + "/" + relPath
}
