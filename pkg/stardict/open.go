package stardict

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func isDir(pathStr string) bool {
	stat, _ := os.Stat(pathStr)
	if stat == nil {
		return false
	}
	return stat.IsDir()
}

// Open open directories
func Open(d string) ([]*Dictionary, error) {
	var dicList []*Dictionary
	const ext = ".ifo"
	filepath.Walk(d, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		name := info.Name()
		if filepath.Ext(info.Name()) != ext {
			return nil
		}
		fmt.Printf("Loading %#v\n", path)
		dirPath := filepath.Dir(path)
		dic, err := NewDictionary(dirPath, name[:len(name)-len(ext)])
		if err != nil {
			fmt.Println(err)
			return err
		}
		resDir := filepath.Join(dirPath, "res")
		if isDir(resDir) {
			dic.resDir = resDir
			dic.resURL = "file://" + pathToUnix(resDir)

		}
		dicList = append(dicList, dic)
		return nil
	})
	return dicList, nil
}

func pathToUnix(pathStr string) string {
	if runtime.GOOS != "windows" {
		return pathStr
	}
	return "/" + strings.Replace(pathStr, `\`, `/`, -1)
}
