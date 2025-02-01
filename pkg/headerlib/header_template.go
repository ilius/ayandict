package headerlib

import (
	"bytes"
	"html"
	"html/template"
	"strings"

	"github.com/ilius/ayandict/v2/pkg/config"
	"github.com/ilius/ayandict/v2/pkg/dictmgr"
	"github.com/ilius/ayandict/v2/pkg/wordwrap"
	commons "github.com/ilius/go-dict-commons"
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
) (string, error) {
	terms := res.Terms()
	termsJoined := html.EscapeString(strings.Join(terms, " | "))
	headerBuf := bytes.NewBuffer(nil)
	dictName := res.DictName()
	err := headerTpl.Execute(headerBuf, HeaderTemplateInput{
		Terms:     terms,
		Term:      termsJoined,
		DictName:  dictName,
		Score:     res.Score() >> 1,
		ShowTerms: dictmgr.DictShowTerms(dictName),
	})
	if err != nil {
		return "", err
	}
	return headerBuf.String(), nil
}
