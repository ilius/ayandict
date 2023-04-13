package stardict

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	I_bookname    = "bookname"
	I_wordcount   = "wordcount"
	I_description = "description"
	I_idxfilesize = "idxfilesize"
)

// Info contains dictionary options
type Info struct {
	Options  map[string]string
	Version  string
	Is64     bool
	disabled bool
}

func (info Info) DictName() string {
	return info.Options[I_bookname]
}

// EntryCount returns number of words in the dictionary
func (info Info) EntryCount() (int, error) {
	num, err := strconv.ParseUint(info.Options[I_wordcount], 10, 64)
	if err != nil {
		return 0, err
	}
	return int(num), nil
}

func (info Info) Description() string {
	return info.Options[I_description]
}

func (info Info) IndexFileSize() uint64 {
	num, _ := strconv.ParseUint(info.Options[I_idxfilesize], 10, 64)
	return num
}

func (info Info) MaxIdxBytes() int {
	if info.Is64 {
		return 8
	}
	return 4
}

func decodeOption(str string) (key string, value string, err error) {
	a := strings.Split(str, "=")

	if len(a) < 2 {
		return "", "", errors.New("Invalid file format: " + str)
	}

	return a[0], a[1], nil
}

// ReadInfo reads ifo file and collects dictionary options
func ReadInfo(filename string) (info *Info, err error) {
	reader, err := os.Open(filename)
	if err != nil {
		return
	}

	defer reader.Close()

	r := bufio.NewReader(reader)

	_, err = r.ReadString('\n')

	if err != nil {
		return
	}

	version, err := r.ReadString('\n')
	if err != nil {
		return
	}

	key, value, err := decodeOption(version[:len(version)-1])
	if err != nil {
		return
	}

	if key != "version" {
		err = errors.New("Version missing (should be on second line)")
		return
	}

	if value != "2.4.2" && value != "3.0.0" {
		err = errors.New("Stardict version should be either 2.4.2 or 3.0.0")
		return
	}

	info = new(Info)

	info.Version = value

	info.Options = make(map[string]string)

	for {
		option, err := r.ReadString('\n')

		if err != nil && err != io.EOF {
			return info, err
		}

		if err == io.EOF && len(option) == 0 {
			break
		}

		key, value, err = decodeOption(option[:len(option)-1])

		if err != nil {
			return info, err
		}

		info.Options[key] = value

		if err == io.EOF {
			break
		}
	}

	if bits, ok := info.Options["idxoffsetbits"]; ok {
		if bits == "64" {
			info.Is64 = true
		}
	} else {
		info.Is64 = false
	}

	return
}
