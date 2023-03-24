package stardict

import (
	"bytes"
	"encoding/binary"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Translation contains translation items
type Translation struct {
	Parts []*TranslationItem
}

// TranslationItem contain single translation item
type TranslationItem struct {
	Type rune
	Data []byte
}

type SearchResult struct {
	Term  string
	Items []*TranslationItem
}

// Dictionary stardict dictionary
type Dictionary struct {
	dict *Dict
	idx  *Idx
	info *Info

	resDir string
	resURL string
}

func (d *Dictionary) ResourceDir() string {
	return d.resDir
}

func (d *Dictionary) ResourceURL() string {
	return d.resURL
}

func (d *Dictionary) translate(offset uint64, size uint64) (items []*TranslationItem) {
	if _, ok := d.info.Options["sametypesequence"]; ok {
		return d.translateWithSametypesequence(d.dict.GetSequence(offset, size))
	}
	return d.translateWithoutSametypesequence(d.dict.GetSequence(offset, size))
}

// Search: first try an exact match
// then search all translations for terms that contain the query
// but sort the one that have it as prefix first
func (d *Dictionary) Search(query string) []*SearchResult {
	// if len(query) < 2 {
	// 	return d.searchVeryShort(query)
	// }
	idx := d.idx
	results0 := []*SearchResult{}
	results1 := []*SearchResult{}
	results2 := []*SearchResult{}

	query = strings.ToLower(strings.TrimSpace(query))
	queryWords := strings.Split(query, " ")

	mainWordIndex := 0
	// we can change this by allowing a query like '* case' to match 'test case' etc

	queryMainWord := queryWords[mainWordIndex]
	prefix := queryMainWord
	if len(queryMainWord) > 2 {
		prefix = queryMainWord[:2]
	}
Loop:
	for _, termIndex := range idx.byWordPrefix[prefix] {
		entry := idx.terms[termIndex]
		term := strings.ToLower(entry.Term)
		if query == term {
			results0 = append(results0, &SearchResult{
				Term:  term,
				Items: d.translate(entry.Offset, entry.Size),
			})
			continue
		}
		termWords := strings.Split(term, " ")
		if len(termWords) <= mainWordIndex {
			continue
		}
		if queryMainWord == termWords[mainWordIndex] {
			results1 = append(results1, &SearchResult{
				Term:  term,
				Items: d.translate(entry.Offset, entry.Size),
			})
			continue
		}
		for _, termWord := range termWords {
			if queryMainWord == termWord {
				results2 = append(results2, &SearchResult{
					Term:  term,
					Items: d.translate(entry.Offset, entry.Size),
				})
				continue Loop
			}
		}
		if strings.Contains(termWords[mainWordIndex], queryMainWord) {
			results2 = append(results2, &SearchResult{
				Term:  term,
				Items: d.translate(entry.Offset, entry.Size),
			})
			continue
		}
	}
	results := results0
	if len(results1) > 0 {
		results = append(results, results1...)
	}
	if len(results2) > 0 {
		results = append(results, results2...)
	}
	return results
}

func (d *Dictionary) translateWithSametypesequence(data []byte) (items []*TranslationItem) {
	seq := d.info.Options["sametypesequence"]

	seqLen := len(seq)

	var dataPos int
	dataSize := len(data)

	for i, t := range seq {
		switch t {
		case 'm', 'l', 'g', 't', 'x', 'y', 'k', 'w', 'h', 'r':
			// if last seq item
			if i == seqLen-1 {
				items = append(items, &TranslationItem{Type: t, Data: data[dataPos:dataSize]})
			} else {
				end := bytes.IndexRune(data[dataPos:], '\000')
				items = append(items, &TranslationItem{Type: t, Data: data[dataPos : dataPos+end+1]})
				dataPos += end + 1
			}
		case 'W', 'P':
			if i == seqLen-1 {
				items = append(items, &TranslationItem{Type: t, Data: data[dataPos:dataSize]})
			} else {
				size := binary.BigEndian.Uint32(data[dataPos : dataPos+4])
				items = append(items, &TranslationItem{Type: t, Data: data[dataPos+4 : dataPos+int(size)+5]})
				dataPos += int(size) + 5
			}
		}
	}

	return
}

func (d *Dictionary) translateWithoutSametypesequence(data []byte) (items []*TranslationItem) {
	var dataPos int
	dataSize := len(data)

	for {
		t := data[dataPos]

		dataPos++

		switch t {
		case 'm', 'l', 'g', 't', 'x', 'y', 'k', 'w', 'h', 'r':
			end := bytes.IndexRune(data[dataPos:], '\000')

			if end < 0 { // last item
				items = append(items, &TranslationItem{Type: rune(t), Data: data[dataPos:dataSize]})
				dataPos = dataSize
			} else {
				items = append(items, &TranslationItem{Type: rune(t), Data: data[dataPos : dataPos+end+1]})
				dataPos += end + 1
			}
		case 'W', 'P':
			size := binary.BigEndian.Uint32(data[dataPos : dataPos+4])
			items = append(items, &TranslationItem{Type: rune(t), Data: data[dataPos+4 : dataPos+int(size)+5]})
			dataPos += int(size) + 5
		}

		if dataPos >= dataSize {
			break
		}
	}

	return
}

// BookName returns book name
func (d *Dictionary) BookName() string {
	return d.info.Options["bookname"]
}

// EntryCount returns number of entries in the dictionary
func (d *Dictionary) EntryCount() uint64 {
	num, _ := strconv.ParseUint(d.info.Options["wordcount"], 10, 64)

	return num
}

// NewDictionary returns a new Dictionary
// path - path to dictionary files
// name - name of dictionary to parse
func NewDictionary(path string, name string) (*Dictionary, error) {
	d := new(Dictionary)

	path = filepath.Clean(path)

	dictDzPath := filepath.Join(path, name+".dict.dz")
	dictPath := filepath.Join(path, name+".dict")

	idxPath := filepath.Join(path, name+".idx")
	infoPath := filepath.Join(path, name+".ifo")

	if _, err := os.Stat(infoPath); os.IsNotExist(err) {
		return nil, err
	}

	if _, err := os.Stat(idxPath); os.IsNotExist(err) {
		return nil, err
	}

	// we should have either .dict.dz or .dict file
	if _, err := os.Stat(dictDzPath); os.IsNotExist(err) {
		if _, err := os.Stat(dictPath); os.IsNotExist(err) {
			return nil, err
		}
	} else {
		dictPath = dictDzPath
	}

	info, err := ReadInfo(infoPath)
	if err != nil {
		return nil, err
	}

	idx, err := ReadIndex(idxPath, info)
	if err != nil {
		return nil, err
	}

	dict, err := ReadDict(dictPath, info)
	if err != nil {
		return nil, err
	}

	d.info = info
	d.idx = idx
	d.dict = dict

	return d, nil
}
