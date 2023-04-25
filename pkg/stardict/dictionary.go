package stardict

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gobwas/glob"
	"github.com/ilius/ayandict/pkg/murmur3"
	common "github.com/ilius/go-dict-commons"
	su "github.com/ilius/go-dict-commons/search_utils"
)

// dictionaryImp stardict dictionary
type dictionaryImp struct {
	*Info

	dict     *Dict
	idx      *Idx
	ifoPath  string
	idxPath  string
	dictPath string
	synPath  string
	resDir   string
	resURL   string

	decodeData func(data []byte) []*common.SearchResultItem
}

func (d *dictionaryImp) Disabled() bool {
	return d.disabled
}

func (d *dictionaryImp) Loaded() bool {
	return d.dict != nil
}

func (d *dictionaryImp) SetDisabled(disabled bool) {
	d.disabled = disabled
}

func (d *dictionaryImp) ResourceDir() string {
	return d.resDir
}

func (d *dictionaryImp) ResourceURL() string {
	return d.resURL
}

func (d *dictionaryImp) IndexPath() string {
	return d.idxPath
}

func (d *dictionaryImp) InfoPath() string {
	return d.ifoPath
}

func (d *dictionaryImp) Close() {
	d.dict.Close()
}

func (d *dictionaryImp) CalcHash() ([]byte, error) {
	file, err := os.Open(d.idxPath)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	hash := murmur3.New128()
	if _, err := io.Copy(hash, file); err != nil {
		return nil, err
	}
	return hash.Sum(nil), nil
}

func (d *dictionaryImp) newResult(entry *IdxEntry, entryIndex int, score uint8) *common.SearchResultLow {
	return &common.SearchResultLow{
		F_Score: score,
		F_Terms: entry.terms,
		Items: func() []*common.SearchResultItem {
			return d.decodeData(d.dict.GetSequence(entry.offset, entry.size))
		},
		F_EntryIndex: uint64(entryIndex),
	}
}

// SearchFuzzy: run a fuzzy search with similarity scores
// ranging from 140 (which means %70) to 200 (which means 100%)
func (d *dictionaryImp) SearchFuzzy(
	query string,
	workerCount int,
	timeout time.Duration,
) []*common.SearchResultLow {
	// if len(query) < 2 {
	// 	return d.searchVeryShort(query)
	// }

	idx := d.idx
	const minScore = uint8(140)

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

	prefix := queryMainWord[0]
	entryIndexes := idx.byWordPrefix[prefix]

	t1 := time.Now()
	N := len(entryIndexes)

	args := &su.ScoreFuzzyArgs{
		Query:          query,
		QueryRunes:     queryRunes,
		QueryMainWord:  queryMainWord,
		QueryWordCount: queryWordCount,
		MinWordCount:   minWordCount,
		MainWordIndex:  mainWordIndex,
	}

	results := su.RunWorkers(
		N,
		workerCount,
		timeout,
		func(start int, end int) []*common.SearchResultLow {
			var results []*common.SearchResultLow
			for i := start; i < end; i++ {
				entry := idx.entries[entryIndexes[i]]
				score := su.ScoreFuzzy(entry.terms, args)
				if score < minScore {
					continue
				}
				results = append(results, d.newResult(entry, i, score))
			}
			return results
		},
	)

	dt := time.Since(t1)
	if dt > time.Millisecond {
		log.Printf("SearchFuzzy index loop took %v for %#v on %s\n", dt, query, d.DictName())
	}
	// log.Printf("Search produced %d results for %#v on %s\n", len(results), query, d.DictName())
	return results
}

func (d *dictionaryImp) SearchStartWith(
	query string,
	workerCount int,
	timeout time.Duration,
) []*common.SearchResultLow {
	idx := d.idx
	const minScore = uint8(140)

	query = strings.ToLower(strings.TrimSpace(query))

	prefix, _ := utf8.DecodeRuneInString(query)
	entryIndexes := idx.byWordPrefix[prefix]

	t1 := time.Now()
	N := len(entryIndexes)

	results := su.RunWorkers(
		N,
		workerCount,
		timeout,
		func(start int, end int) []*common.SearchResultLow {
			var results []*common.SearchResultLow
			for i := start; i < end; i++ {
				entry := idx.entries[entryIndexes[i]]
				score := su.ScoreStartsWith(entry.terms, query)
				if score < minScore {
					continue
				}
				results = append(results, d.newResult(entry, i, score))
			}
			return results
		},
	)

	dt := time.Since(t1)
	if dt > time.Millisecond {
		log.Printf("SearchStartWith index loop took %v for %#v on %s\n", dt, query, d.DictName())
	}
	// log.Printf("Search produced %d results for %#v on %s\n", len(results), query, d.DictName())
	return results
}

func (d *dictionaryImp) searchPattern(
	workerCount int,
	timeout time.Duration,
	checkTerm func(string) uint8,
) []*common.SearchResultLow {
	idx := d.idx
	const minScore = uint8(140)

	N := len(idx.entries)
	return su.RunWorkers(
		N,
		workerCount,
		timeout,
		func(start int, end int) []*common.SearchResultLow {
			var results []*common.SearchResultLow
			for entryI := start; entryI < end; entryI++ {
				entry := idx.entries[entryI]
				score := uint8(0)
				for _, term := range entry.terms {
					termScore := checkTerm(term)
					if termScore > score {
						score = termScore
						break
					}
				}
				if score < minScore {
					continue
				}
				results = append(results, d.newResult(entry, entryI, score))
			}
			return results
		},
	)
}

func (d *dictionaryImp) SearchRegex(
	query string,
	workerCount int,
	timeout time.Duration,
) ([]*common.SearchResultLow, error) {
	re, err := regexp.Compile("^" + query + "$")
	if err != nil {
		return nil, err
	}

	t1 := time.Now()
	results := d.searchPattern(workerCount, timeout, func(term string) uint8 {
		if !re.MatchString(term) {
			return 0
		}
		if len(term) < 20 {
			return 200 - uint8(len(term))
		}
		return 180
	})
	dt := time.Since(t1)
	if dt > time.Millisecond {
		log.Printf("SearchRegex index loop took %v for %#v on %s\n", dt, query, d.DictName())
	}
	return results, nil
}

func (d *dictionaryImp) SearchGlob(
	query string,
	workerCount int,
	timeout time.Duration,
) ([]*common.SearchResultLow, error) {
	pattern, err := glob.Compile(query)
	if err != nil {
		return nil, err
	}

	t1 := time.Now()
	results := d.searchPattern(workerCount, timeout, func(term string) uint8 {
		if !pattern.Match(term) {
			return 0
		}
		if len(term) < 20 {
			return 200 - uint8(len(term))
		}
		return 180
	})
	dt := time.Since(t1)
	if dt > time.Millisecond {
		log.Printf("SearchGlob index loop took %v for %#v on %s\n", dt, query, d.DictName())
	}
	return results, nil
}

func (d *dictionaryImp) decodeWithSametypesequence(data []byte) (items []*common.SearchResultItem) {
	seq := d.Options["sametypesequence"]

	seqLen := len(seq)

	var dataPos int
	dataSize := len(data)

	for i, t := range seq {
		switch t {
		case 'm', 'l', 'g', 't', 'x', 'y', 'k', 'w', 'h', 'r':
			// if last seq item
			if i == seqLen-1 {
				items = append(items, &common.SearchResultItem{Type: t, Data: data[dataPos:dataSize]})
			} else {
				end := bytes.IndexRune(data[dataPos:], '\000')
				items = append(items, &common.SearchResultItem{Type: t, Data: data[dataPos : dataPos+end+1]})
				dataPos += end + 1
			}
		case 'W', 'P':
			if i == seqLen-1 {
				items = append(items, &common.SearchResultItem{Type: t, Data: data[dataPos:dataSize]})
			} else {
				size := binary.BigEndian.Uint32(data[dataPos : dataPos+4])
				items = append(items, &common.SearchResultItem{Type: t, Data: data[dataPos+4 : dataPos+int(size)+5]})
				dataPos += int(size) + 5
			}
		}
	}

	return
}

func (d *dictionaryImp) decodeWithoutSametypesequence(data []byte) (items []*common.SearchResultItem) {
	var dataPos int
	dataSize := len(data)

	for {
		t := data[dataPos]

		dataPos++

		switch t {
		case 'm', 'l', 'g', 't', 'x', 'y', 'k', 'w', 'h', 'r':
			end := bytes.IndexRune(data[dataPos:], '\000')

			if end < 0 { // last item
				items = append(items, &common.SearchResultItem{Type: rune(t), Data: data[dataPos:dataSize]})
				dataPos = dataSize
			} else {
				items = append(items, &common.SearchResultItem{Type: rune(t), Data: data[dataPos : dataPos+end+1]})
				dataPos += end + 1
			}
		case 'W', 'P':
			size := binary.BigEndian.Uint32(data[dataPos : dataPos+4])
			items = append(items, &common.SearchResultItem{Type: rune(t), Data: data[dataPos+4 : dataPos+int(size)+5]})
			dataPos += int(size) + 5
		}

		if dataPos >= dataSize {
			break
		}
	}

	return
}

// DictName returns book name
func (d *dictionaryImp) DictName() string {
	return d.Options["bookname"]
}

// NewDictionary returns a new Dictionary
// path - path to dictionary files
// name - name of dictionary to parse
func NewDictionary(path string, name string) (*dictionaryImp, error) {
	d := &dictionaryImp{}

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

	// we should have either .dict or .dict.dz file
	if _, err := os.Stat(dictPath); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		if _, errDz := os.Stat(dictDzPath); errDz != nil {
			return nil, err
		}
		dictPath = dictDzPath
	}

	info, err := ReadInfo(ifoPath)
	if err != nil {
		return nil, err
	}
	d.Info = info

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

func (d *dictionaryImp) Load() error {
	{
		idx, err := ReadIndex(d.idxPath, d.synPath, d.Info)
		if err != nil {
			return err
		}
		d.idx = idx
	}
	{
		dict, err := ReadDict(d.dictPath, d.Info)
		if err != nil {
			return err
		}
		d.dict = dict
	}
	return nil
}
