package stardict

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/ilius/ayandict/pkg/levenshtein"
)

// SearchResultItem contain single translation item
type SearchResultItem struct {
	Data []byte
	Type rune
}

type SearchResult struct {
	terms []string
	items func() []*SearchResultItem
	score uint8
}

// Dictionary stardict dictionary
type Dictionary struct {
	disabled bool

	ifoPath  string
	idxPath  string
	dictPath string
	synPath  string

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

func (d *Dictionary) translate(offset uint64, size uint64) []*SearchResultItem {
	if _, ok := d.info.Options["sametypesequence"]; ok {
		return d.translateWithSametypesequence(d.dict.GetSequence(offset, size))
	}
	return d.translateWithoutSametypesequence(d.dict.GetSequence(offset, size))
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
func (d *Dictionary) Search(query string, cutoff int) []*SearchResult {
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

	chechEntry := func(entry *IdxEntry) uint8 {
		bestScore := uint8(0)
		for _, termOrig := range entry.Terms {
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
			if bestWordScore > 50 {
				if queryWordCount > 1 {
					bestWordScore = bestWordScore/2 + bestWordScore/7
				}
				if bestWordScore > score {
					score = bestWordScore
				}
			}
			if score > bestScore {
				bestScore = score
			}
		}
		return bestScore
	}

	prefix, _ := utf8.DecodeRuneInString(queryMainWord)
	for _, termIndex := range idx.byWordPrefix[prefix] {
		entry := idx.terms[termIndex]
		score := chechEntry(entry)
		if score > 100 {
			results = append(results, &SearchResult{
				score: score,
				terms: entry.Terms,
				items: func() []*SearchResultItem {
					return d.translate(entry.Offset, entry.Size)
				},
			})
		}
	}
	fmt.Printf("Search produced %d results for %#v on %s\n", len(results), query, d.DictName())
	sort.Slice(results, func(i, j int) bool {
		return results[i].score > results[j].score
	})
	if cutoff > 0 && len(results) > cutoff {
		for ; cutoff < len(results) && results[cutoff].score > 180; cutoff++ {
		}
		results = results[:cutoff]
	}
	return results
}

func (d *Dictionary) translateWithSametypesequence(data []byte) (items []*SearchResultItem) {
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

func (d *Dictionary) translateWithoutSametypesequence(data []byte) (items []*SearchResultItem) {
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
