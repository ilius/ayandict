package dictmgr

import (
	"fmt"
	std_html "html"

	"github.com/ilius/ayandict/pkg/config"
	common "github.com/ilius/go-dict-commons"
)

type SearchResult struct {
	*common.SearchResultLow
	dic    common.Dictionary
	conf   *config.Config
	hDefis []string
}

func (r *SearchResult) DictName() string {
	return r.dic.DictName()
}

func (r *SearchResult) DefinitionsHTML() []string {
	if r.hDefis != nil {
		return r.hDefis
	}
	definitions := []string{}
	resURL := r.dic.ResourceURL()
	for _, item := range r.Items() {
		if item.Type == 'h' {
			itemDefi := string(item.Data)
			itemDefi = fixDefiHTML(itemDefi, resURL, r.conf, r.dic)
			definitions = append(definitions, itemDefi+"<br/>\n")
			continue
		}
		definitions = append(definitions, fmt.Sprintf(
			"<pre>%s</pre>\n<br/>\n",
			std_html.EscapeString(string(item.Data)),
		))
	}
	r.hDefis = definitions
	return definitions
}

func (r *SearchResult) ResourceDir() string {
	return r.dic.ResourceDir()
}
