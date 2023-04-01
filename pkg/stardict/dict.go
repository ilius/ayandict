package stardict

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
)

type t_ReadSeekerCloser interface {
	io.ReadSeeker
	io.Closer
}

// Dict implements in-memory dictionary
type Dict struct {
	filename string

	r    t_ReadSeekerCloser
	lock sync.Mutex
}

// GetSequence returns data at the given offset
func (d *Dict) GetSequence(offset uint64, size uint64) []byte {
	if d.r == nil {
		log.Println("GetSequence: file is closed")
		return nil
	}
	d.lock.Lock()
	defer d.lock.Unlock()
	d.r.Seek(int64(offset), 0)
	p := make([]byte, size)
	_, err := d.r.Read(p)
	if err != nil {
		log.Printf("error while reading dict file %#v: %v\n", d.filename, err)
		return nil
	}
	return p
}

func (d *Dict) Open() error {
	reader, err := os.Open(d.filename)
	if err != nil {
		return err
	}
	d.r = reader
	return nil
}

func (d *Dict) Close() {
	if d.r == nil {
		return
	}
	fmt.Println("Closing", d.filename)
	d.r.Close()
	d.r = nil
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
