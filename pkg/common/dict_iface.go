package common

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
	CalcHash() ([]byte, error)
	SearchFuzzy(query string) []*SearchResultLow
}