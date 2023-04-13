package common

type DictSettings struct {
	Symbol string `json:"symbol"`
	Order  int    `json:"order"`
	Hash   string `json:"hash"`
}

func NewDictSettings(dic Dictionary, index int) *DictSettings {
	return &DictSettings{
		Symbol: DefaultSymbol(dic.DictName()),
		Order:  index,
		Hash:   Hash(dic),
	}
}
