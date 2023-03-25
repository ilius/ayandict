package common

type QueryResult struct {
	Term        string
	DictName    string
	Definitions []string
	Score       uint8
}
