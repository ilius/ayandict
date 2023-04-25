package application

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/ilius/ayandict/pkg/config"
	"github.com/ilius/ayandict/pkg/qerr"
)

var (
	history          = []string{}
	historyMaxSize   = 100
	historyMutex     sync.Mutex
	historySaveMutex sync.Mutex
	frequencyMutex   sync.Mutex
)

const (
	historyFileName   = "history.json"
	frequencyFileName = "frequent.json"
)

func addHistoryLow(query string) {
	historyMutex.Lock()
	history = append(history, query)
	if len(history) > historyMaxSize {
		history = history[len(history)-historyMaxSize:]
	}
	historyMutex.Unlock()
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
	jsonBytes, err := os.ReadFile(pathStr)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("error loading history: %v", err)
		}
		return nil
	}
	err = json.Unmarshal(jsonBytes, &history)
	if err != nil {
		return fmt.Errorf("bad history file %#v: %v", pathStr, err)
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
	err = os.WriteFile(historyFilePath(), jsonBytes, 0o644)
	if err != nil {
		qerr.Errorf("Error saving history: %v", err)
	}
}
