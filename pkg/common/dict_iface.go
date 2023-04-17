package common

import "time"

type Dictionary interface {
	Disabled() bool
	SetDisabled(bool)
	Loaded() bool
	Load() error
	Close()
	DictName() string
	EntryCount() (int, error)
	Description() string
	ResourceDir() string
	ResourceURL() string
	IndexPath() string
	IndexFileSize() uint64
	InfoPath() string
	CalcHash() ([]byte, error)
	SearchFuzzy(query string, workerCount int, timeout time.Duration) []*SearchResultLow
	SearchStartWith(query string, workerCount int, timeout time.Duration) []*SearchResultLow
	SearchRegex(query string, workerCount int, timeout time.Duration) ([]*SearchResultLow, error)
	SearchGlob(query string, workerCount int, timeout time.Duration) ([]*SearchResultLow, error)
}
