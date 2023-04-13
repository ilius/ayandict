package common

type SearchResultIface interface {
	Terms() []string
	Score() uint8
	DictName() string
	DefinitionsHTML() []string
	ResourceDir() string
}

type SearchResultItem struct {
	Data []byte
	Type rune
}

type SearchResultLow struct {
	Items   func() []*SearchResultItem
	F_Terms []string
	F_Score uint8
}

func (res *SearchResultLow) Terms() []string {
	return res.F_Terms
}

func (res *SearchResultLow) Score() uint8 {
	return res.F_Score
}
