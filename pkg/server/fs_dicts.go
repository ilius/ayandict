package server

import (
	"embed"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/ilius/ayandict/v2/pkg/dictmgr"
)

type dictResFileSystem struct {
	fs embed.FS
	// prefix string
}

func (f *dictResFileSystem) Open(name string) (http.File, error) {
	name = strings.TrimLeft(name, "/")
	parts := strings.Split(name, "/")
	if len(parts) < 3 {
		return nil, os.ErrNotExist
	}
	dictName := parts[1]
	resPath := filepath.Join(parts[2:]...)
	// fmt.Printf("dictName=%#v, resPath=%#v", dictName, resPath)
	fpath, ok := dictmgr.DictResFile(dictName, resPath)
	if !ok {
		return nil, os.ErrNotExist
	}
	file, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}
	return &httpFile{file}, nil
}
