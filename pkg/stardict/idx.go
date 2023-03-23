package stardict

import (
	"encoding/binary"
	"io/ioutil"
	"strconv"
)

// Sense has information belonging to single item position in dictionary
type Sense = [2]uint64

// Idx implements an in-memory index for a dictionary
type Idx struct {
	items map[string][]Sense
}

// NewIdx initializes idx struct
func NewIdx(wordCount int) *Idx {
	idx := new(Idx)
	if wordCount > 0 {
		idx.items = make(map[string][]Sense, wordCount)
	} else {
		idx.items = make(map[string][]Sense)
	}
	return idx
}

// Add adds an item to in-memory index
func (idx *Idx) Add(item string, offset uint64, size uint64) {
	idx.items[item] = append(idx.items[item], Sense{offset, size})
}

// ReadIndex reads dictionary index into a memory and returns in-memory index structure
func ReadIndex(filename string, info *Info) (idx *Idx, err error) {
	data, err := ioutil.ReadFile(filename)
	// unable to read index
	if err != nil {
		return
	}

	wordCount := 0
	wordCountStr := info.Options["wordcount"]
	if wordCountStr != "" {
		n, err := strconv.ParseInt(wordCountStr, 10, 64)
		if err != nil {
			return nil, err
		}
		wordCount = int(n)
	}

	idx = NewIdx(wordCount)

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
