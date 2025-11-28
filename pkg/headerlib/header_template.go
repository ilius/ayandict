package headerlib

import (
	"bytes"
	"html"
	"html/template"

	commons "codeberg.org/ilius/go-dict-commons"
	"github.com/ilius/ayandict/v3/pkg/config"
	"github.com/ilius/ayandict/v3/pkg/dictmgr"
	"github.com/ilius/ayandict/v3/pkg/wordwrap"
)

type HeaderTemplateInput struct {
	Terms     []string
	Term      string
	DictName  string
	Score     uint8
	ShowTerms bool
}

func LoadHeaderTemplate(conf *config.Config) (*template.Template, error) {
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

func GetHeader(
	headerTpl *template.Template,
	res commons.SearchResultIface,
	maxTermsTextLength int,
) (string, error) {
	terms := res.Terms()
	terms, termsJoined := joinWithMaxLen(terms, " | ", maxTermsTextLength)
	headerBuf := bytes.NewBuffer(nil)
	dictName := res.DictName()
	err := headerTpl.Execute(headerBuf, HeaderTemplateInput{
		Terms:     terms,
		Term:      html.EscapeString(termsJoined),
		DictName:  dictName,
		Score:     res.Score() >> 1,
		ShowTerms: dictmgr.DictShowTerms(dictName),
	})
	if err != nil {
		return "", err
	}
	return headerBuf.String(), nil
}
