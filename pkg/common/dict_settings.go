package common

type DictSettings struct {
	Symbol string `json:"symbol"`
	Order  int    `json:"order"`
	Hash   string `json:"hash"`
}

func NewDictSettings(info Dictionary, index int) *DictSettings {
	return &DictSettings{
		Symbol: DefaultSymbol(info.DictName()),
		Order:  index,
		Hash:   Hash(info),
	}
}
