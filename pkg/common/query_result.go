package common

type QueryResult struct {
	Score       uint8
	Term        string
	DictName    string
	Definitions []string
}
