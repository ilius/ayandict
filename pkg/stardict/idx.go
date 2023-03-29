package stardict

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"unicode/utf8"
)

type IdxEntry struct {
	Terms  []string
	Offset uint64
	Size   uint64
}

// Idx implements an in-memory index for a dictionary
type Idx struct {
	byWordPrefix map[rune][]int
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
	idx.byWordPrefix = map[rune][]int{}
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

// ReadIndex reads dictionary index into a memory and returns in-memory index structure
func ReadIndex(filename string, synPath string, info *Info) (*Idx, error) {
	data, err := ioutil.ReadFile(filename)
	// unable to read index
	if err != nil {
		return nil, err
	}

	entryCount := 0
	entryCountStr := info.Options["wordcount"]
	if entryCountStr != "" {
		n, err := strconv.ParseInt(entryCountStr, 10, 64)
		if err != nil {
			return nil, err
		}
		entryCount = int(n)
	}

	idx := NewIdx(entryCount)

	var a [255]byte // temporary buffer
	var aIdx int
	var expect int

	var term string
	var dataOffset uint64
	var dataSize uint64

	maxIntBytes := info.MaxIdxBytes()

	byWordPrefix := map[rune]map[int]bool{}

	addTermPrefix := func(term string, termIndex int) {
		for _, word := range strings.Split(strings.ToLower(term), " ") {
			prefix, _ := utf8.DecodeRuneInString(word)
			m, ok := byWordPrefix[prefix]
			if !ok {
				m = map[int]bool{}
				byWordPrefix[prefix] = m
			}
			m[termIndex] = true
		}
	}

	for _, b := range data {
		if expect == 0 {
			a[aIdx] = b
			if b == 0 {
				term = string(a[:aIdx])

				aIdx = 0
				expect++
				continue
			}
			aIdx++
			continue
		}
		if expect == 1 {
			a[aIdx] = b
			if aIdx == maxIntBytes-1 {
				if info.Is64 {
					dataOffset = binary.BigEndian.Uint64(a[:maxIntBytes])
				} else {
					dataOffset = uint64(binary.BigEndian.Uint32(a[:maxIntBytes]))
				}

				aIdx = 0
				expect++
				continue
			}
			aIdx++
			continue
		}
		a[aIdx] = b
		if aIdx == maxIntBytes-1 {
			if info.Is64 {
				dataSize = binary.BigEndian.Uint64(a[:maxIntBytes])
			} else {
				dataSize = uint64(binary.BigEndian.Uint32(a[:maxIntBytes]))
			}

			aIdx = 0
			expect = 0

			// finished with one record
			termIndex := idx.Add(term, dataOffset, dataSize)

			addTermPrefix(term, termIndex)

			continue
		}
		aIdx++
	}
	if synPath != "" {
		data, err := ioutil.ReadFile(synPath)
		// unable to read index
		if err != nil {
			return nil, err
		}
		dataLen := len(data)
		pos := 0
		for pos < dataLen {
			beg := pos
			// Python: pos = data.find("\x00", beg)
			offset := bytes.Index(data[beg:], []byte{0})
			if offset < 0 {
				return nil, fmt.Errorf("Synonym file is corrupted")
			}
			pos = offset + beg
			b_alt := data[beg:pos]
			pos += 1
			if pos+4 > len(data) {
				return nil, fmt.Errorf("Synonym file is corrupted")
			}
			termIndex := int(binary.BigEndian.Uint32(data[pos : pos+4]))
			pos += 4
			if termIndex >= len(idx.terms) {
				return nil, fmt.Errorf(
					"Corrupted synonym file. Word %#v references invalid item",
					string(b_alt),
				)
			}
			alt := string(b_alt)
			entry := idx.terms[termIndex]
			entry.Terms = append(entry.Terms, alt)
			addTermPrefix(alt, termIndex)
		}
	}

	for prefix, indexMap := range byWordPrefix {
		indexList := make([]int, 0, len(indexMap))
		for i := range indexMap {
			indexList = append(indexList, i)
		}
		idx.byWordPrefix[prefix] = indexList
	}

	return idx, err
}
