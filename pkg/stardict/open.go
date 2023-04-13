package stardict

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/ilius/ayandict/pkg/qerr"
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
		log.Printf("Initializing %#v\n", path)
		dirPath := filepath.Dir(path)
		dic, err := NewDictionary(dirPath, name[:len(name)-len(ext)])
		if err != nil {
			return err
		}
		if order[dic.DictName()] < 0 {
			dic.Disabled = true
			dicList = append(dicList, dic)
			return nil
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
			log.Println(err)
			qerr.Error(err)
		}
	}
	var wg sync.WaitGroup
	load := func(dic *Dictionary) {
		defer wg.Done()
		// log.Printf("Loading index %#v\n", dic.idxPath)
		err = dic.load()
		if err != nil {
			qerr.Errorf("error loading %#v: %v", dic.DictName(), err)
			log.Printf("error loading %#v: %v", dic.DictName(), err)
		} else {
			log.Printf("Loaded index %#v\n", dic.idxPath)
		}
	}
	for _, dic := range dicList {
		if dic.Disabled {
			continue
		}
		wg.Add(1)
		go load(dic)
	}
	wg.Wait()
	return dicList, nil
}

func pathToUnix(pathStr string) string {
	if runtime.GOOS != "windows" {
		return pathStr
	}
	return "/" + strings.Replace(pathStr, `\`, `/`, -1)
}
