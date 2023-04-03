package stardict

import (
	"compress/gzip"
	"io"
	"log"
	"os"
	"strings"
	"sync"
)

// Dict implements in-memory dictionary
type Dict struct {
	filename string

	file *os.File
	lock sync.Mutex
}

func (d *Dict) Open() error {
	file, err := os.Open(d.filename)
	if err != nil {
		return err
	}
	d.file = file
	return nil
}

func (d *Dict) Close() {
	if d.file == nil {
		return
	}
	log.Println("Closing", d.filename)
	d.file.Close()
	d.file = nil
}

func dictunzip(filename string) (string, error) {
	newFilename := filename[0 : len(filename)-3]
	reader, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	gzReader, err := gzip.NewReader(reader)
	defer gzReader.Close()
	writer, err := os.Create(newFilename)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(writer, gzReader)
	if err != nil {
		return "", err
	}
	return newFilename, nil
}

// ReadDict creates Dict and opens .dict file
func ReadDict(filename string, info *Info) (*Dict, error) {
	if strings.HasSuffix(filename, ".dz") {
		// if file is compressed then read it from archive
		var err error
		filename, err = dictunzip(filename)
		if err != nil {
			return nil, err
		}
	}
	dict := &Dict{
		filename: filename,
	}
	err := dict.Open()
	if err != nil {
		return nil, err
	}
	return dict, nil
}
