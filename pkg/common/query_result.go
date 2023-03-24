package common

type QueryResult struct {
	Score       float64
	Term        string
	DictName    string
	Definitions []string
}
