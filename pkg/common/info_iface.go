package common

type Info interface {
	DictName() string
	EntryCount() (int, error)
	Description() string
	IndexFileSize() uint64
	CalcHash() ([]byte, error)
}
