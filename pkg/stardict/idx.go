package stardict

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"

	"github.com/dolthub/swiss"
)

type IdxEntry struct {
	Terms  []string
	Offset uint64
	Size   uint64
}

// Idx implements an in-memory index for a dictionary
type Idx struct {
	byWordPrefix *swiss.Map[rune, []int]
	terms        []*IdxEntry
}

// NewIdx initializes idx struct
func NewIdx(entryCount int) *Idx {
	idx := new(Idx)
	if entryCount > 0 {
		idx.terms = make([]*IdxEntry, 0, entryCount)
	} else {
		idx.terms = []*IdxEntry{}
	}
	idx.byWordPrefix = swiss.NewMap[rune, []int](100)
	return idx
}

// Add adds an item to in-memory index
func (idx *Idx) Add(term string, offset uint64, size uint64) int {
	termIndex := len(idx.terms)
	idx.terms = append(idx.terms, &IdxEntry{
		Terms:  []string{term},
		Offset: offset,
		Size:   size,
	})
	return termIndex
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
	idx := NewIdx(entryCount)

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
		termIndex := idx.Add(term, dataOffset, num)
		wordPrefixMap.Add(term, termIndex)
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
		idx.byWordPrefix.Put(prefix, indexList)
	}

	return idx, err
}

func readSyn(idx *Idx, synPath string, wordPrefixMap WordPrefixMap) error {
	data, err := ioutil.ReadFile(synPath)
	// unable to read index
	if err != nil {
		return err
	}
	dataLen := len(data)
	pos := 0
	for pos < dataLen {
		beg := pos
		// Python: pos = data.find("\x00", beg)
		offset := bytes.Index(data[beg:], []byte{0})
		if offset < 0 {
			return fmt.Errorf("Synonym file is corrupted")
		}
		pos = offset + beg
		b_alt := data[beg:pos]
		pos += 1
		if pos+4 > len(data) {
			return fmt.Errorf("Synonym file is corrupted")
		}
		termIndex := int(binary.BigEndian.Uint32(data[pos : pos+4]))
		pos += 4
		if termIndex >= len(idx.terms) {
			return fmt.Errorf(
				"Corrupted synonym file. Word %#v references invalid item",
				string(b_alt),
			)
		}
		alt := string(b_alt)
		entry := idx.terms[termIndex]
		entry.Terms = append(entry.Terms, alt)
		wordPrefixMap.Add(alt, termIndex)
	}
	return nil
}
