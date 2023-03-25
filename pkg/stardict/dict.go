package stardict

import (
	"compress/gzip"
	"io"
	"log"
	"os"
	"strings"
)

type t_ReadSeekerCloser interface {
	io.ReadSeeker
	io.Closer
}

// Dict implements in-memory dictionary
type Dict struct {
	r        t_ReadSeekerCloser
	filename string
}

// GetSequence returns data at the given offset
func (d Dict) GetSequence(offset uint64, size uint64) []byte {
	d.r.Seek(int64(offset), 0)
	p := make([]byte, size)
	_, err := d.r.Read(p)
	if err != nil {
		log.Printf("error while reading dict file %#v: %v\n", d.filename, err)
		return nil
	}
	return p
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

// ReadDict reads dictionary into memory
func ReadDict(filename string, info *Info) (dict *Dict, err error) {
	if strings.HasSuffix(filename, ".dz") { // if file is compressed then read it from archive
		filename, err = dictunzip(filename)
		if err != nil {
			return
		}
	}

	reader, err := os.Open(filename)
	if err != nil {
		return
	}

	dict = new(Dict)
	dict.filename = filename
	dict.r = reader

	return
}
