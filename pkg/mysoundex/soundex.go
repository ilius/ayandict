package mysoundex

import (
	"bufio"
	"log/slog"
	"os"
	"sync"

	"github.com/xrash/smetrics"
)

func NewSoundexSearcher(wordsFile string) *SoundexSearcher {
	ss := &SoundexSearcher{
		wordsFile: wordsFile,

		m: map[string][]string{},
	}

	return ss
}

type SoundexSearcher struct {
	wordsFile string

	m    map[string][]string
	lock sync.RWMutex
}

func (ss *SoundexSearcher) Load() {
	wordsFile := ss.wordsFile
	ss.lock.Lock()
	defer ss.lock.Unlock()
	file, err := os.Open(wordsFile)
	if err != nil {
		slog.Error("failed to open soundex words file", "err", err, "path", wordsFile)
		return
	}
	scanner := bufio.NewScanner(file)
	slog.Info("----------- reading soundex words file", "path", wordsFile)
	for scanner.Scan() {
		word := scanner.Text()
		if word == "" {
			continue
		}
		code := smetrics.Soundex(word)
		if code == "0000" {
			continue
		}
		ss.m[code] = append(ss.m[code], word)
	}
	if err := scanner.Err(); err != nil {
		slog.Error("error reading file", "err", err, "path", wordsFile)
		return
	}
	slog.Info("--------- finished reading soundex words file", "path", wordsFile)
}

func (ss *SoundexSearcher) Lookup(word string) []string {
	if !ss.lock.TryRLock() {
		slog.Error("soundex words file is still loading")
		return nil
	}
	defer ss.lock.RUnlock()
	return ss.m[smetrics.Soundex(word)]
}
