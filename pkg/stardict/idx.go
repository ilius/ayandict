package stardict

import (
	"encoding/binary"
	"io/ioutil"
	"strings"
)

type IdxEntry struct {
	terms  []string
	offset uint64
	size   uint64
}

// Idx implements an in-memory index for a dictionary
type Idx struct {
	byWordPrefix map[rune][]int
	entries      []*IdxEntry
}

// newIdx initializes idx struct
func newIdx(entryCount int) *Idx {
	idx := &Idx{
		byWordPrefix: map[rune][]int{},
	}
	if entryCount > 0 {
		idx.entries = make([]*IdxEntry, 0, entryCount)
	} else {
		idx.entries = []*IdxEntry{}
	}
	return idx
}

type t_state int

const (
	termState t_state = iota
	offsetState
	sizeState
)

// ReadIndex reads dictionary index into a memory and returns in-memory index structure
func ReadIndex(filename string, synPath string, info *Info) (*Idx, error) {
	data, err := ioutil.ReadFile(filename)
	// unable to read index
	if err != nil {
		return nil, err
	}

	entryCount, err := info.EntryCount()
	if err != nil {
		return nil, err
	}
	idx := newIdx(entryCount)

	wordPrefixMap := WordPrefixMap{}

	var buf [255]byte // temporary buffer
	var bufPos int
	state := termState

	var term string
	var dataOffset uint64

	maxIntBytes := info.MaxIdxBytes()

	for _, b := range data {
		buf[bufPos] = b
		if state == termState {
			if b > 0 {
				bufPos++
				continue
			}
			term = string(buf[:bufPos])
			bufPos = 0
			state = offsetState
			continue
		}
		if bufPos < maxIntBytes-1 {
			bufPos++
			continue
		}
		var num uint64
		if info.Is64 {
			num = binary.BigEndian.Uint64(buf[:maxIntBytes])
		} else {
			num = uint64(binary.BigEndian.Uint32(buf[:maxIntBytes]))
		}
		if state == offsetState {
			dataOffset = num
			bufPos = 0
			state = sizeState
			continue
		}
		// finished with one record
		bufPos = 0
		state = termState
		words := strings.Split(strings.ToLower(term), " ")
		termIndex := len(idx.entries)
		idx.entries = append(idx.entries, &IdxEntry{
			terms:  []string{term},
			offset: dataOffset,
			size:   num,
		})
		wordPrefixMap.Add(words, termIndex)
	}
	if synPath != "" {
		err := readSyn(idx, synPath, wordPrefixMap)
		if err != nil {
			return nil, err
		}
	}
	for prefix, indexMap := range wordPrefixMap {
		indexList := make([]int, 0, len(indexMap))
		for i := range indexMap {
			indexList = append(indexList, i)
		}
		idx.byWordPrefix[prefix] = indexList
	}

	return idx, err
}
