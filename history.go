package main

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
	history          = []string{}
	historyMutex     sync.Mutex
	historySaveMutex sync.Mutex
)

const historyFileName = "history.json"

var addHistoryGUI func(string)

func addHistoryLow(query string) {
	historyMutex.Lock()
	history = append(history, query)
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
}

func historyFilePath() string {
	return filepath.Join(config.GetConfigDir(), historyFileName)
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
	pathStr := historyFilePath()
	jsonBytes, err := json.MarshalIndent(history, "", "\t")
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(pathStr, jsonBytes, 0o644)
	if err != nil {
		fmt.Printf("Error saving history: %v\n", err)
	}
}
