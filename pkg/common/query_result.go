package common

type QueryResult interface {
	Terms() []string
	DictName() string
	DefinitionsHTML() []string
	Score() uint8
}
