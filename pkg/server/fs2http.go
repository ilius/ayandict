package server

import (
	"embed"
	"errors"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
)

type httpFileSystem struct {
	fs     embed.FS
	prefix string
}

type localFile interface {
	fs.File
	io.Seeker
}

type httpFile struct {
	localFile
}

func (f *httpFile) Readdir(count int) ([]fs.FileInfo, error) {
	return nil, errors.New("not a directory")
}

func (f *httpFileSystem) Open(name string) (http.File, error) {
	name = strings.TrimLeft(name, "/")
	if name == f.prefix {
		name = f.prefix + "/index.html"
	}
	file, err := f.fs.Open(name)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	file2, ok := file.(localFile)
	if !ok {
		log.Printf("file %#v is not a seeker", name)
		return nil, os.ErrNotExist
	}
	return &httpFile{file2}, nil
}
