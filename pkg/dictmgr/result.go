package dictmgr

import (
	"fmt"
	std_html "html"

	"github.com/ilius/ayandict/v2/pkg/config"
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

func (r *SearchResult) ResourceDir() string {
	return r.dic.ResourceDir()
}

func (r *SearchResult) DefinitionsHTML(flags uint32) []string {
	if r.hDefis != nil {
		return r.hDefis
	}
	definitions := []string{}
	for _, item := range r.Items() {
		if item.Type == 'h' {
			itemDefi := string(item.Data)
			itemDefi = fixDefiHTML(
				itemDefi,
				r.conf,
				r.dic,
				flags,
			)
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
