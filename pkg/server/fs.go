package server

import (
	"errors"
	"io"
	"io/fs"
)

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
