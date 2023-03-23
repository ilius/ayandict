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
	Keyword string
	Items   []*TranslationItem
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

func (d *Dictionary) translate(senses [][2]uint64) (items []*Translation) {
	for _, sense := range senses {
		data := d.dict.GetSequence(sense[0], sense[1])

		var transItems []*TranslationItem

		if _, ok := d.info.Options["sametypesequence"]; ok {
			transItems = d.translateWithSametypesequence(data)
		} else {
			transItems = d.translateWithoutSametypesequence(data)
		}

		items = append(items, &Translation{Parts: transItems})
	}

	return
}

func (d *Dictionary) searchVeryShort(query string) []*SearchResult {
	terms := []string{query}
	queryLower := strings.ToLower(query)
	if queryLower != query {
		terms = append(terms, queryLower)
	}
	queryUpper := strings.ToUpper(query)
	if queryUpper != query {
		terms = append(terms, queryUpper)
	}
	results := []*SearchResult{}
	for _, term := range terms {
		senses := d.idx.items[term]
		if senses == nil {
			continue
		}
		result := &SearchResult{
			Keyword: term,
		}
		for _, item := range d.translate(senses) {
			result.Items = append(result.Items, item.Parts...)
		}
		results = append(results, result)
	}
	return results
}

// SearchAuto: first try an exact match
// then search all translations for keywords that contain the query
// but sort the one that have it as prefix first
func (d *Dictionary) SearchAuto(query string) []*SearchResult {
	if len(query) < 2 {
		return d.searchVeryShort(query)
	}
	results0 := []*SearchResult{}
	results1 := []*SearchResult{}
	results2 := []*SearchResult{}
	for keyword, senses := range d.idx.items {
		if keyword == query {
			result := &SearchResult{Keyword: keyword}
			for _, item := range d.translate(senses) {
				result.Items = append(result.Items, item.Parts...)
			}
			results0 = append(results0, result)
			continue
		}
		if strings.HasPrefix(keyword, query) {
			result := &SearchResult{Keyword: keyword}
			for _, item := range d.translate(senses) {
				result.Items = append(result.Items, item.Parts...)
			}
			results1 = append(results1, result)
			continue
		}
		if strings.Contains(keyword, query) {
			result := &SearchResult{Keyword: keyword}
			for _, item := range d.translate(senses) {
				result.Items = append(result.Items, item.Parts...)
			}
			results2 = append(results2, result)
		}
	}
	results := append(results0, results1...)
	return append(results, results2...)
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

// GetBookName returns book name
func (d *Dictionary) GetBookName() string {
	return d.info.Options["bookname"]
}

// GetWordCount returns number of words in the dictionary
func (d *Dictionary) GetWordCount() uint64 {
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
