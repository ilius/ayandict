package stardict

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
)

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
		if termIndex >= len(idx.entries) {
			return fmt.Errorf(
				"Corrupted synonym file. Word %#v references invalid item",
				string(b_alt),
			)
		}
		alt := []rune(string(b_alt))
		entry := idx.entries[termIndex]
		entry.terms = append(entry.terms, alt)
		wordPrefixMap.Add(alt, termIndex)
	}
	return nil
}
