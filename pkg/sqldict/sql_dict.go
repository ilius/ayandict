package sqldict

import (
	"database/sql"
	"database/sql/driver"
	"encoding/hex"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	// _ "github.com/glebarez/go-sqlite"
	"github.com/ilius/ayandict/pkg/qerr"
	su "github.com/ilius/ayandict/pkg/search_utils"
	common "github.com/ilius/go-dict-commons"
	"modernc.org/sqlite"
)

const minScore = uint8(140)

// uriList[i] == "sqlite://PATH.db"
func Open(uriList []string, order map[string]int) []common.Dictionary {
	dicList := []common.Dictionary{}
	for _, uri := range uriList {
		i := strings.Index(uri, "://")
		if i < 1 {
			qerr.Errorf("invalid sql dict uri = %#v", uri)
			continue
		}
		driver := uri[:i]
		source := uri[i+3:]
		dic := NewDictionary(driver, source)
		name := dic.DictName()
		if order[name] < 0 {
			dic.disabled = true
		}
		dicList = append(dicList, dic)
	}
	return dicList
}

func NewDictionary(driver string, source string) *dictionaryImp {
	return &dictionaryImp{
		driver: driver,
		source: source,
	}
}

type dictionaryImp struct {
	disabled bool
	dictName string
	driver   string
	source   string
	hash     []byte

	db *sql.DB
}

func (d *dictionaryImp) Disabled() bool {
	return d.disabled
}

func (d *dictionaryImp) SetDisabled(disabled bool) {
	d.disabled = disabled
}

func (d *dictionaryImp) Loaded() bool {
	return d.db != nil
}

func (d *dictionaryImp) Load() error {
	err := d.defineRegexp()
	if err != nil {
		return err
	}
	db, err := sql.Open(d.driver, d.source)
	if err != nil {
		return err
	}
	d.db = db
	return nil
}

func (d *dictionaryImp) Close() {
	if d.db == nil {
		return
	}
	err := d.db.Close()
	if err != nil {
		log.Println(err)
	}
	d.db = nil
}

func (d *dictionaryImp) readInfo(key string) (string, error) {
	if d.db == nil {
		err := d.Load()
		if err != nil {
			return "", err
		}
	}
	row := d.db.QueryRow("SELECT value FROM meta WHERE key = ?", key)
	if row == nil {
		return "", fmt.Errorf("no %v in meta", key)
	}
	value := ""
	err := row.Scan(&value)
	if err != nil {
		return "", err
	}
	log.Printf("info: %v = %#v", key, value)
	return value, nil
}

func (d *dictionaryImp) DictName() string {
	if d.dictName != "" {
		return d.dictName
	}
	dictName, err := d.readInfo("name")
	if err != nil {
		log.Println(err)
		d.dictName = d.source
		return d.dictName
	}
	d.dictName = dictName
	return dictName
}

func (d *dictionaryImp) EntryCount() (int, error) {
	row := d.db.QueryRow("SELECT count(id) FROM entry")
	if row == nil {
		return 0, fmt.Errorf("EntryCount: row = nil")
	}
	count := 0
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (d *dictionaryImp) Description() string {
	desc, err := d.readInfo("description")
	if err != nil {
		log.Println(err)
		return ""
	}
	return desc
}

func (d *dictionaryImp) ResourceDir() string {
	// TODO
	return ""
}

func (d *dictionaryImp) ResourceURL() string {
	// TODO
	return ""
}

func (d *dictionaryImp) IndexPath() string {
	return ""
}

func (d *dictionaryImp) IndexFileSize() uint64 {
	return 0
}

func (d *dictionaryImp) InfoPath() string {
	return ""
}

func (d *dictionaryImp) CalcHash() ([]byte, error) {
	if d.hash != nil {
		return d.hash, nil
	}
	hexStr, err := d.readInfo("hash")
	if err != nil {
		return nil, err
	}
	hash := make([]byte, len(hexStr)/2)
	_, err = hex.Decode(hash, []byte(hexStr))
	if err != nil {
		return nil, err
	}
	d.hash = hash
	return hash, nil
}

func (d *dictionaryImp) readArticle(id int) []*common.SearchResultItem {
	row := d.db.QueryRow("SELECT article FROM entry WHERE id IS ?", id)
	if row == nil {
		log.Printf("No row with id = %v", id)
		return nil
	}
	article := ""
	err := row.Scan(&article)
	if err != nil {
		log.Println(err)
		return nil
	}
	return []*common.SearchResultItem{
		{
			Type: 'h',
			Data: []byte(article),
		},
	}
}

func (d *dictionaryImp) newResult(terms []string, id int, score uint8) *common.SearchResultLow {
	return &common.SearchResultLow{
		F_Score: score,
		F_Terms: terms,
		Items: func() []*common.SearchResultItem {
			return d.readArticle(id)
		},
		F_EntryIndex: uint64(id),
	}
}

func (d *dictionaryImp) SearchFuzzy(query string, _ int, _ time.Duration) []*common.SearchResultLow {
	t1 := time.Now()
	query = strings.ToLower(strings.TrimSpace(query))

	rows, err := d.db.Query(
		"SELECT id, term FROM entry WHERE term LIKE ?",
		"%"+query+"%",
	)
	if err != nil {
		qerr.Error(err)
		return nil
	}
	queryWords := strings.Split(query, " ")

	mainWordIndex := 0
	for mainWordIndex < len(queryWords)-1 && queryWords[mainWordIndex] == "*" {
		mainWordIndex++
	}

	minWordCount := 1
	queryWordCount := 0
	for _, word := range queryWords {
		if word == "*" {
			minWordCount++
			continue
		}
		queryWordCount++
	}

	args := &su.ScoreFuzzyArgs{
		Query:          query,
		QueryRunes:     []rune(query),
		QueryMainWord:  []rune(queryWords[mainWordIndex]),
		QueryWordCount: queryWordCount,
		MinWordCount:   minWordCount,
		MainWordIndex:  mainWordIndex,
	}
	results := []*common.SearchResultLow{}
	for rows.Next() {
		id := -1
		term := ""
		err := rows.Scan(&id, &term)
		if err != nil {
			qerr.Error(err)
			return nil
		}
		// TODO: alts
		terms := []string{term}
		score := su.ScoreEntryFuzzy(terms, args)
		if score < minScore {
			continue
		}
		results = append(results, d.newResult(terms, id, score))
	}
	dt := time.Since(t1)
	if dt > time.Millisecond {
		log.Printf("SearchFuzzy index loop took %v for %#v on %s\n", dt, query, d.DictName())
	}
	return results
}

func (d *dictionaryImp) SearchStartWith(query string, _ int, _ time.Duration) []*common.SearchResultLow {
	t1 := time.Now()
	query = strings.ToLower(strings.TrimSpace(query))
	rows, err := d.db.Query(
		"SELECT id, term FROM entry WHERE term LIKE ?",
		query+"%",
	)
	if err != nil {
		qerr.Error(err)
		return nil
	}
	results := []*common.SearchResultLow{}
	for rows.Next() {
		id := -1
		term := ""
		err := rows.Scan(&id, &term)
		if err != nil {
			qerr.Error(err)
			return nil
		}
		// TODO: alts
		terms := []string{term}
		score := su.ScoreStartsWith(terms, query)
		if score < minScore {
			continue
		}
		results = append(results, d.newResult(terms, id, score))
	}
	dt := time.Since(t1)
	if dt > time.Millisecond {
		log.Printf("SearchStartWith index loop took %v for %#v on %s\n", dt, query, d.DictName())
	}
	return results
}

func (d *dictionaryImp) searchPattern(
	where string,
	arg string,
	checkTerm func(string) uint8,
) []*common.SearchResultLow {
	sqlQ := "SELECT id, term FROM entry WHERE " + where
	rows, err := d.db.Query(sqlQ, arg)
	if err != nil {
		qerr.Errorf("error running SQL query %#v: %v", sqlQ, err)
		return nil
	}
	results := []*common.SearchResultLow{}
	for rows.Next() {
		id := -1
		term := ""
		err := rows.Scan(&id, &term)
		if err != nil {
			qerr.Error(err)
			return nil
		}
		// TODO: alts
		terms := []string{term}
		score := uint8(0)
		for _, term := range terms {
			termScore := checkTerm(term)
			if termScore > score {
				score = termScore
				break
			}
		}
		if score < minScore {
			continue
		}
		results = append(results, d.newResult(terms, id, score))
	}
	return results
}

func (d *dictionaryImp) SearchRegex(query string, _ int, _ time.Duration) ([]*common.SearchResultLow, error) {
	t1 := time.Now()
	results := d.searchPattern("term REGEXP ?", "^"+query+"$", func(term string) uint8 {
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

func (d *dictionaryImp) SearchGlob(query string, _ int, _ time.Duration) ([]*common.SearchResultLow, error) {
	t1 := time.Now()
	results := d.searchPattern("term GLOB ?", query, func(term string) uint8 {
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

func (d *dictionaryImp) defineRegexp() error {
	const argc = 2
	return sqlite.RegisterDeterministicScalarFunction(
		"regexp",
		argc,
		func(ctx *sqlite.FunctionContext, args []driver.Value) (driver.Value, error) {
			s1 := args[0].(string)
			s2 := args[1].(string)
			matched, err := regexp.MatchString(s1, s2)
			if err != nil {
				return nil, fmt.Errorf("bad regular expression: %q", err)
			}
			// sqlite3.Xsqlite3_result_int(tls, ctx, libc.Bool32(matched))
			return matched, nil
		},
	)
}
