package dictmgr

import (
	"os"
	"sync"
)

var (
	imageTempMap      = map[string]string{}
	imageTempMapMutex sync.RWMutex
)

func loadPNGFile(filename string) (string, error) {
	imageTempMapMutex.RLock()
	tmpPath, ok := imageTempMap[filename]
	imageTempMapMutex.RUnlock()
	if ok {
		return tmpPath, nil
	}
	imageTempMapMutex.Lock()
	defer imageTempMapMutex.Unlock()
	data, err := res.ReadFile("res/" + filename)
	if err != nil {
		return "", err
	}
	file, err := os.CreateTemp("", filename)
	if err != nil {
		return "", err
	}
	_, err = file.Write(data)
	if err != nil {
		return "", err
	}
	err = file.Close()
	if err != nil {
		return "", err
	}
	imageTempMap[filename] = file.Name()
	return file.Name(), nil
}
