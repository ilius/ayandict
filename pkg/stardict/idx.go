package stardict

import (
	"encoding/binary"
	"io/ioutil"
	"strconv"
	"strings"
)

type IdxEntry struct {
	Term   string
	Offset uint64
	Size   uint64
}

// Idx implements an in-memory index for a dictionary
type Idx struct {
	terms []IdxEntry

	byWordPrefix map[string][]int
}

// NewIdx initializes idx struct
func NewIdx(entryCount int) *Idx {
	idx := new(Idx)
	if entryCount > 0 {
		idx.terms = make([]IdxEntry, entryCount)
	} else {
		idx.terms = []IdxEntry{}
	}
	idx.byWordPrefix = map[string][]int{}
	return idx
}

// Add adds an item to in-memory index
func (idx *Idx) Add(term string, offset uint64, size uint64) {
	termIndex := len(idx.terms)
	idx.terms = append(idx.terms, IdxEntry{
		Term:   term,
		Offset: offset,
		Size:   size,
	})
	for _, word := range strings.Split(strings.ToLower(term), " ") {
		var prefix string
		if len(word) > 2 {
			prefix = word[:2]
		} else {
			prefix = word
		}
		idx.byWordPrefix[prefix] = append(idx.byWordPrefix[prefix], termIndex)
	}
}

// ReadIndex reads dictionary index into a memory and returns in-memory index structure
func ReadIndex(filename string, info *Info) (idx *Idx, err error) {
	data, err := ioutil.ReadFile(filename)
	// unable to read index
	if err != nil {
		return
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

	idx = NewIdx(entryCount)

	var a [255]byte // temporary buffer
	var aIdx int
	var expect int

	var dataStr string
	var dataOffset uint64
	var dataSize uint64

	var maxIntBytes int

	if info.Is64 == true {
		maxIntBytes = 8
	} else {
		maxIntBytes = 4
	}

	for _, b := range data {
		if expect == 0 {
			a[aIdx] = b
			if b == 0 {
				dataStr = string(a[:aIdx])

				aIdx = 0
				expect++
				continue
			}
			aIdx++
		} else {
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
			} else {
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
					idx.Add(dataStr, dataOffset, dataSize)

					continue
				}
				aIdx++
			}
		}
	}

	return idx, err
}
