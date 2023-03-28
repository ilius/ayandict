package application

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/ilius/ayandict/pkg/config"
)

var (
	history            = []string{}
	historyMaxSize     = 100
	historyMutex       sync.Mutex
	historySaveMutex   sync.Mutex
	frequencyMutex     sync.Mutex
	frequencySaveMutex sync.Mutex
)

const (
	historyFileName   = "history.json"
	frequencyFileName = "frequent.json"
)

var addHistoryGUI func(string)

var trimHistoryGUI func(int)

func addHistoryAndFrequency(query string) {
	if !conf.HistoryDisable {
		addHistory(query)
	}
	if !conf.MostFrequentDisable {
		frequencyView.Add(query, 1)
		if conf.MostFrequentAutoSave {
			SaveFrequency()
		}
	}
}

func addHistoryLow(query string) {
	historyMutex.Lock()
	history = append(history, query)
	if len(history) > historyMaxSize {
		history = history[len(history)-historyMaxSize:]
	}
	historyMutex.Unlock()
}

func addHistory(query string) {
	if len(history) > 0 && query == history[len(history)-1] {
		return
	}
	addHistoryLow(query)
	if addHistoryGUI != nil {
		addHistoryGUI(query)
	}
	if trimHistoryGUI != nil {
		trimHistoryGUI(historyMaxSize)
	}
	if conf.HistoryAutoSave {
		SaveHistory()
	}
}

func historyFilePath() string {
	return filepath.Join(config.GetConfigDir(), historyFileName)
}

func frequencyFilePath() string {
	return filepath.Join(config.GetConfigDir(), frequencyFileName)
}

func LoadHistory() error {
	historyMutex.Lock()
	defer historyMutex.Unlock()
	pathStr := historyFilePath()
	jsonBytes, err := ioutil.ReadFile(pathStr)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("Error loading history: %v\n", err)
		}
		return nil
	}
	err = json.Unmarshal(jsonBytes, &history)
	if err != nil {
		return fmt.Errorf("Bad history file %#v: %v\n", pathStr, err)
	}
	return nil
}

func SaveHistory() {
	historySaveMutex.Lock()
	defer historySaveMutex.Unlock()
	jsonBytes, err := json.MarshalIndent(history, "", "\t")
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(historyFilePath(), jsonBytes, 0o644)
	if err != nil {
		fmt.Printf("Error saving history: %v\n", err)
	}
}

func SaveFrequency() {
	frequencySaveMutex.Lock()
	defer frequencySaveMutex.Unlock()
	err := frequencyView.SaveToFile(frequencyFilePath())
	if err != nil {
		fmt.Printf("Error saving history: %v\n", err)
	}
}

func clearHistory() {
	historyMutex.Lock()
	history = []string{}
	historyMutex.Unlock()

	frequencyView.Clear()

	SaveHistory()
	SaveFrequency()
}
