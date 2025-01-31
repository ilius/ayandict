package config

import (
	"html/template"

	"github.com/ilius/ayandict/v2/pkg/wordwrap"
)

type HeaderTemplateInput struct {
	Terms     []string
	Term      string
	DictName  string
	Score     uint8
	ShowTerms bool
}

func LoadHeaderTemplate(conf *Config) (*template.Template, error) {
	// slog.Info("Parsing:", conf.HeaderTemplate)
	tpl := template.New("header").Funcs(template.FuncMap{
		"wrapterms": func(terms []string, limit int) [][]string {
			return wordwrap.WordWrapByWords(terms, limit, " ", " ")
		},
	})
	tpl, err := tpl.Parse(conf.HeaderTemplate)
	if err != nil {
		return nil, err
	}
	return tpl, nil
}
