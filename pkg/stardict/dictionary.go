package stardict

import (
	"bytes"
	"encoding/binary"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/ilius/ayandict/pkg/levenshtein"
)

// SearchResultItem contain single translation item
type SearchResultItem struct {
	Data []byte
	Type rune
}

type SearchResult struct {
	items func() []*SearchResultItem
	terms []string
	score uint8
}

// Dictionary stardict dictionary
type Dictionary struct {
	dict     *Dict
	idx      *Idx
	info     *Info
	ifoPath  string
	idxPath  string
	dictPath string
	synPath  string
	resDir   string
	resURL   string
	disabled bool

	decodeData func(data []byte) []*SearchResultItem
}

func (d *Dictionary) ResourceDir() string {
	return d.resDir
}

func (d *Dictionary) ResourceURL() string {
	return d.resURL
}

func similarity(a string, b string) uint8 {
	n := len(a)
	if len(b) > n {
		n = len(b)
	}
	return uint8(200 * (n - levenshtein.ComputeDistance(a, b)) / n)
}

// Search: first try an exact match
// then search all translations for terms that contain the query
// but sort the one that have it as prefix first
func (d *Dictionary) Search(query string) []*SearchResult {
	// if len(query) < 2 {
	// 	return d.searchVeryShort(query)
	// }
	idx := d.idx
	results := []*SearchResult{}

	query = strings.ToLower(strings.TrimSpace(query))
	queryWords := strings.Split(query, " ")

	mainWordIndex := 0
	for mainWordIndex < len(queryWords)-1 && queryWords[mainWordIndex] == "*" {
		mainWordIndex++
	}
	queryMainWord := queryWords[mainWordIndex]

	minWordCount := 1
	queryWordCount := 0
	for _, word := range queryWords {
		if word == "*" {
			minWordCount++
			continue
		}
		queryWordCount++
	}

	chechEntry := func(terms []string) uint8 {
		bestScore := uint8(0)
		for _, termOrig := range terms {
			term := strings.ToLower(termOrig)
			if term == query {
				return 200
			}
			if strings.Contains(term, query) {
				score := uint8(200 * (1 + len(query)) / (1 + len(term)))
				if score > bestScore {
					bestScore = score
					continue
				}
			}
			words := strings.Split(term, " ")
			if len(words) < minWordCount {
				continue
			}
			score := similarity(query, term)
			if score > bestScore {
				bestScore = score
				continue
			}
			// if score < 50 {
			// 	continue
			// }
			if len(words) > 1 {
				bestWordScore := uint8(0)
				for wordI, word := range words {
					wordScore := similarity(queryMainWord, word)
					if wordI != mainWordIndex {
						wordScore -= wordScore / 10
					}
					if wordScore < 140 {
						continue
					}
					if wordScore > bestWordScore {
						bestWordScore = wordScore
					}
				}
				if bestWordScore < 50 {
					continue
				}
				if queryWordCount > 1 {
					bestWordScore = bestWordScore/2 + bestWordScore/7
				}
				if bestWordScore > bestScore {
					bestScore = bestWordScore
				}
			}
		}
		return bestScore
	}

	t1 := time.Now()
	prefix, _ := utf8.DecodeRuneInString(queryMainWord)

	const minScore = uint8(140)

	for _, termIndex := range idx.byWordPrefix[prefix] {
		entry := idx.terms[termIndex]
		score := chechEntry(entry.Terms)
		if score < minScore {
			continue
		}
		results = append(results, &SearchResult{
			score: score,
			terms: entry.Terms,
			items: func() []*SearchResultItem {
				return d.decodeData(d.dict.GetSequence(entry.Offset, entry.Size))
			},
		})

	}
	dt := time.Now().Sub(t1)
	if dt > time.Millisecond {
		log.Printf("Search index loop took %v for %#v on %s\n", dt, query, d.DictName())
	}
	// log.Printf("Search produced %d results for %#v on %s\n", len(results), query, d.DictName())
	return results
}

func (d *Dictionary) decodeWithSametypesequence(data []byte) (items []*SearchResultItem) {
	seq := d.info.Options["sametypesequence"]

	seqLen := len(seq)

	var dataPos int
	dataSize := len(data)

	for i, t := range seq {
		switch t {
		case 'm', 'l', 'g', 't', 'x', 'y', 'k', 'w', 'h', 'r':
			// if last seq item
			if i == seqLen-1 {
				items = append(items, &SearchResultItem{Type: t, Data: data[dataPos:dataSize]})
			} else {
				end := bytes.IndexRune(data[dataPos:], '\000')
				items = append(items, &SearchResultItem{Type: t, Data: data[dataPos : dataPos+end+1]})
				dataPos += end + 1
			}
		case 'W', 'P':
			if i == seqLen-1 {
				items = append(items, &SearchResultItem{Type: t, Data: data[dataPos:dataSize]})
			} else {
				size := binary.BigEndian.Uint32(data[dataPos : dataPos+4])
				items = append(items, &SearchResultItem{Type: t, Data: data[dataPos+4 : dataPos+int(size)+5]})
				dataPos += int(size) + 5
			}
		}
	}

	return
}

func (d *Dictionary) decodeWithoutSametypesequence(data []byte) (items []*SearchResultItem) {
	var dataPos int
	dataSize := len(data)

	for {
		t := data[dataPos]

		dataPos++

		switch t {
		case 'm', 'l', 'g', 't', 'x', 'y', 'k', 'w', 'h', 'r':
			end := bytes.IndexRune(data[dataPos:], '\000')

			if end < 0 { // last item
				items = append(items, &SearchResultItem{Type: rune(t), Data: data[dataPos:dataSize]})
				dataPos = dataSize
			} else {
				items = append(items, &SearchResultItem{Type: rune(t), Data: data[dataPos : dataPos+end+1]})
				dataPos += end + 1
			}
		case 'W', 'P':
			size := binary.BigEndian.Uint32(data[dataPos : dataPos+4])
			items = append(items, &SearchResultItem{Type: rune(t), Data: data[dataPos+4 : dataPos+int(size)+5]})
			dataPos += int(size) + 5
		}

		if dataPos >= dataSize {
			break
		}
	}

	return
}

// DictName returns book name
func (d *Dictionary) DictName() string {
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

	ifoPath := filepath.Join(path, name+".ifo")
	idxPath := filepath.Join(path, name+".idx")
	synPath := filepath.Join(path, name+".syn")

	dictDzPath := filepath.Join(path, name+".dict.dz")
	dictPath := filepath.Join(path, name+".dict")

	if _, err := os.Stat(ifoPath); err != nil {
		return nil, err
	}
	if _, err := os.Stat(idxPath); err != nil {
		return nil, err
	}
	if _, err := os.Stat(synPath); err != nil {
		synPath = ""
	}

	// we should have either .dict.dz or .dict file
	if _, err := os.Stat(dictDzPath); os.IsNotExist(err) {
		if _, err := os.Stat(dictPath); os.IsNotExist(err) {
			return nil, err
		}
	} else {
		dictPath = dictDzPath
	}

	info, err := ReadInfo(ifoPath)
	if err != nil {
		return nil, err
	}
	d.info = info

	d.ifoPath = ifoPath
	d.idxPath = idxPath
	d.synPath = synPath
	d.dictPath = dictPath

	if _, ok := info.Options["sametypesequence"]; ok {
		d.decodeData = d.decodeWithSametypesequence
	} else {
		d.decodeData = d.decodeWithoutSametypesequence
	}

	return d, nil
}

func (d *Dictionary) load() error {
	{
		idx, err := ReadIndex(d.idxPath, d.synPath, d.info)
		if err != nil {
			return err
		}
		d.idx = idx
	}
	{
		dict, err := ReadDict(d.dictPath, d.info)
		if err != nil {
			return err
		}
		d.dict = dict
	}
	return nil
}
