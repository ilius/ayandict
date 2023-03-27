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
func Open(dirPathList []string, order map[string]int) ([]*Dictionary, error) {
	var dicList []*Dictionary
	const ext = ".ifo"

	walkFunc := func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.IsDir() {
			return nil
		}
		name := fi.Name()
		if filepath.Ext(fi.Name()) != ext {
			return nil
		}
		fmt.Printf("Loading %#v\n", path)
		dirPath := filepath.Dir(path)
		dic, err := NewDictionary(dirPath, name[:len(name)-len(ext)])
		if err != nil {
			return err
		}
		if order[dic.BookName()] < 0 {
			dic.disabled = true
			dic.info.Disabled = true
			dicList = append(dicList, dic)
			return nil
		}
		err = dic.load()
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
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	for _, dirPath := range dirPathList {
		// dirPath = pathFromUnix(dirPath) // not needed for relative paths
		if !filepath.IsAbs(dirPath) {
			dirPath = filepath.Join(homeDir, dirPath)
		}
		err := filepath.Walk(dirPath, walkFunc)
		if err != nil {
			fmt.Println(err)
		}
	}
	return dicList, nil
}

func pathToUnix(pathStr string) string {
	if runtime.GOOS != "windows" {
		return pathStr
	}
	return "/" + strings.Replace(pathStr, `\`, `/`, -1)
}
