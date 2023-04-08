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

func similarity(r1 []rune, r2 []rune) uint8 {
	if len(r1) > len(r2) {
		r1, r2 = r2, r1
	}
	// now len(r1) <= len(r2)
	n := len(r2)
	if len(r1) < n*2/3 {
		// this optimization assumes we want to ignore below %66 similarity
		return 0
	}
	return uint8(200 * (n - levenshtein.ComputeDistance(r1, r2)) / n)
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
	queryRunes := []rune(query)

	mainWordIndex := 0
	for mainWordIndex < len(queryWords)-1 && queryWords[mainWordIndex] == "*" {
		mainWordIndex++
	}
	queryMainWord := []rune(queryWords[mainWordIndex])

	minWordCount := 1
	queryWordCount := 0
	for _, word := range queryWords {
		if word == "*" {
			minWordCount++
			continue
		}
		queryWordCount++
	}

	chechEntry := func(entry *IdxEntry) uint8 {
		terms := entry.Terms
		bestScore := uint8(0)
		for _, termOrig := range terms {
			term := strings.ToLower(termOrig)
			if term == query {
				return 200
			}
			words := strings.Split(term, " ")
			if len(words) < minWordCount {
				continue
			}
			score := similarity(queryRunes, []rune(term))
			if score > bestScore {
				bestScore = score
				if score >= 180 {
					continue
				}
			}
			if len(words) > 1 {
				bestWordScore := uint8(0)
				for wordI, word := range words {
					wordScore := similarity(queryMainWord, []rune(word))
					if wordScore < 50 {
						continue
					}
					if wordI == mainWordIndex {
						wordScore -= 1
					} else {
						wordScore -= wordScore / 10
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
	prefix := queryMainWord[0]

	const minScore = uint8(140)

	for _, termIndex := range idx.byWordPrefix[prefix] {
		entry := idx.terms[termIndex]
		score := chechEntry(entry)
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
