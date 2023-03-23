package stardict

import (
	"encoding/binary"
	"io/ioutil"
	"strconv"
)

type Senses interface {
	Word() string
	Next() (bool, uint64, uint64)
}

type Senses64 struct {
	word  string
	index int
	items [][2]uint64
}

func (ss *Senses64) Word() string {
	return ss.word
}

func (ss *Senses64) Next() (bool, uint64, uint64) {
	if ss.index >= len(ss.items) {
		return false, 0, 0
	}
	item := ss.items[ss.index]
	ss.index += 1
	return true, item[0], item[1]
}

// Idx implements an in-memory index for a dictionary
type Idx struct {
	items map[string][][2]uint64
}

// NewIdx initializes idx struct
func NewIdx(wordCount int) *Idx {
	idx := new(Idx)
	if wordCount > 0 {
		idx.items = make(map[string][][2]uint64, wordCount)
	} else {
		idx.items = make(map[string][][2]uint64)
	}
	return idx
}

func (idx *Idx) Get(word string) Senses {
	return &Senses64{
		word:  word,
		items: idx.items[word],
	}
}

func (idx *Idx) ForEach(run func(Senses)) {
	for word, rawSenses := range idx.items {
		run(&Senses64{
			word:  word,
			items: rawSenses,
		})
	}
}

// Add adds an item to in-memory index
func (idx *Idx) Add(item string, offset uint64, size uint64) {
	idx.items[item] = append(idx.items[item], [2]uint64{offset, size})
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
