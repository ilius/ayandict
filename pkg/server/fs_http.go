package server

import (
	"embed"
	"net/http"
	"os"
	"strings"
)

type httpFileSystem struct {
	fs     embed.FS
	prefix string
}

func (f *httpFileSystem) Open(name string) (http.File, error) {
	name = strings.TrimLeft(name, "/")
	if name == f.prefix {
		name = f.prefix + "/index.html"
	}
	file, err := f.fs.Open(name)
	if err != nil {
		logger.Error("error opening file", "err", err, "name", name)
		return nil, err
	}
	file2, ok := file.(localFile)
	if !ok {
		logger.Error("file is not a seeker", "name", name)
		return nil, os.ErrNotExist
	}
	return &httpFile{file2}, nil
}
