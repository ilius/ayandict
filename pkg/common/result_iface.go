package common

type QueryResult interface {
	Terms() []string
	DictName() string
	Score() uint8
	DefinitionsHTML() []string
	ResourceDir() string
}
