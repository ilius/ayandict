package stardict

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/ilius/ayandict/pkg/qerr"
	common "github.com/ilius/go-dict-commons"
)

func isDir(pathStr string) bool {
	stat, _ := os.Stat(pathStr)
	if stat == nil {
		return false
	}
	return stat.IsDir()
}

// Open open directories
func Open(dirPathList []string, order map[string]int) ([]common.Dictionary, error) {
	var dicList []common.Dictionary
	const ext = ".ifo"

	findIfoFile := func(path string) (string, os.FileInfo, error) {
		dirEntries, err := os.ReadDir(path)
		if err != nil {
			return "", nil, err
		}
		for _, de := range dirEntries {
			if filepath.Ext(de.Name()) != ext {
				continue
			}
			fi, err := de.Info()
			if err != nil {
				return "", nil, err
			}
			if fi == nil {
				return "", nil, nil
			}
			return filepath.Join(path, fi.Name()), fi, nil
		}
		return "", nil, nil
	}

	checkDirEntry := func(path string, fi os.FileInfo) error {
		if fi.IsDir() {
			ifoPath, ifoFi, err := findIfoFile(path)
			if err != nil {
				return err
			}
			if ifoFi == nil {
				return nil
			}
			fi = ifoFi
			path = ifoPath
		}
		name := fi.Name()
		if filepath.Ext(name) != ext {
			return nil
		}
		log.Printf("Initializing %#v\n", name)
		dirPath := filepath.Dir(path)
		dic, err := NewDictionary(dirPath, name[:len(name)-len(ext)])
		if err != nil {
			return err
		}
		if order[dic.DictName()] < 0 {
			dic.disabled = true
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

		dirEntries, err := ioutil.ReadDir(dirPath)
		if err != nil {
			qerr.Error(err)
			continue
		}
		for _, fi := range dirEntries {
			err := checkDirEntry(filepath.Join(dirPath, fi.Name()), fi)
			if err != nil {
				go qerr.Error(err)
			}
		}
	}
	log.Println("Starting to load indexes")
	var wg sync.WaitGroup
	load := func(dic common.Dictionary) {
		defer wg.Done()
		err = dic.Load()
		if err != nil {
			qerr.Errorf("error loading %#v: %v", dic.DictName(), err)
			log.Printf("error loading %#v: %v", dic.DictName(), err)
		} else {
			log.Printf("Loaded index %#v\n", dic.IndexPath())
		}
	}
	for _, dic := range dicList {
		if dic.Disabled() {
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
