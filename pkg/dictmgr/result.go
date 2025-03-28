package dictmgr

import (
	"fmt"
	std_html "html"

	"github.com/ilius/ayandict/v3/pkg/config"
	common "github.com/ilius/go-dict-commons"
)

func NewSearchResult(
	res *common.SearchResultLow,
	dic common.Dictionary,
	conf *config.Config,
	flags uint32,
) *SearchResult {
	return &SearchResult{
		SearchResultLow: res,

		proc: NewDictProcessor(dic, conf, flags),
	}
}

type SearchResult struct {
	*common.SearchResultLow
	proc   *DictProcessor
	hDefis []string
}

func (r *SearchResult) DictName() string {
	return r.proc.DictName()
}

func (r *SearchResult) ResourceDir() string {
	return r.proc.ResourceDir()
}

func (r *SearchResult) DefinitionsHTML() []string {
	if r.hDefis != nil {
		return r.hDefis
	}
	definitions := []string{}
	for _, item := range r.Items() {
		if item.Type == 'h' {
			itemDefi := string(item.Data)
			itemDefi = r.proc.FixDefiHTML(itemDefi)
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
