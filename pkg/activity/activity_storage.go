package activity

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/ilius/ayandict/v3/pkg/config"
)

func NewActivityStorage(conf *config.Config, configDir string) *ActivityStorage {
	return &ActivityStorage{
		frequencyFilePath: filepath.Join(configDir, "frequent.json"),
		historyFilePath:   filepath.Join(configDir, "history.json"),
		frequencyMap:      map[string]int{},
		saveMutex:         &sync.Mutex{},
		historyMutex:      &sync.Mutex{},
		historyMaxSize:    conf.HistoryMaxSize,
	}
}

type ActivityStorage struct {
	frequencyFilePath string
	historyFilePath   string

	frequencyMap map[string]int

	saveMutex *sync.Mutex

	history        []string
	historyMutex   *sync.Mutex
	historyMaxSize int
}

type FrequencyItem struct {
	Word  string
	Count int
}

func (s *ActivityStorage) AddFrequency(word string, plus int) {
	s.frequencyMap[word] += plus
}

func (s *ActivityStorage) ClearFrequency() {
	s.frequencyMap = map[string]int{}
}

func (s *ActivityStorage) LoadFrequency() ([]*FrequencyItem, error) {
	pathStr := s.frequencyFilePath
	jsonBytes, err := os.ReadFile(pathStr)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("error loading frequency: %w", err)
		}
		return nil, nil
	}
	countMap := map[string]int{}
	err = json.Unmarshal(jsonBytes, &countMap)
	if err != nil {
		return nil, fmt.Errorf("bad frequency file %#v: %w", pathStr, err)
	}
	s.frequencyMap = countMap
	countList := []*FrequencyItem{}
	for key, count := range countMap {
		countList = append(countList, &FrequencyItem{Word: key, Count: count})
	}
	sort.Slice(countList, func(i, j int) bool {
		return countList[i].Count > countList[j].Count
	})
	return countList, nil
}

func (s *ActivityStorage) SaveFrequency() error {
	s.saveMutex.Lock()
	defer s.saveMutex.Unlock()
	jsonBytes, err := json.MarshalIndent(s.frequencyMap, "", "\t")
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(s.frequencyFilePath, jsonBytes, 0o644)
	if err != nil {
		return err
	}
	return nil
}

func (s *ActivityStorage) LoadHistory() ([]string, error) {
	s.historyMutex.Lock()
	defer s.historyMutex.Unlock()
	pathStr := s.historyFilePath
	jsonBytes, err := os.ReadFile(pathStr)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("error loading history: %w", err)
		}
		return nil, nil
	}
	err = json.Unmarshal(jsonBytes, &s.history)
	if err != nil {
		return nil, fmt.Errorf("bad history file %#v: %w", pathStr, err)
	}
	return s.history, nil
}

func (h *ActivityStorage) SaveHistory() error {
	h.saveMutex.Lock()
	defer h.saveMutex.Unlock()
	jsonBytes, err := json.MarshalIndent(h.history, "", "\t")
	if err != nil {
		return err
	}
	return os.WriteFile(h.historyFilePath, jsonBytes, 0o644)
}

func (h *ActivityStorage) AddHistory(query string) bool {
	if len(h.history) > 0 && query == h.history[len(h.history)-1] {
		return false
	}
	h.historyMutex.Lock()
	h.history = append(h.history, query)
	if len(h.history) > h.historyMaxSize {
		h.history = h.history[len(h.history)-h.historyMaxSize:]
	}
	h.historyMutex.Unlock()
	return true
}

func (h *ActivityStorage) ClearHistory() {
	h.historyMutex.Lock()
	h.history = []string{}
	h.historyMutex.Unlock()
}
