package main

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/ilius/ayandict/pkg/config"
)

var commentMap = map[string]string{
	"directory_list": "List of dictionary directory paths",
	"style":          "Path to application stylesheet file (.qss)",
	"article_style":  "Path to article stylesheet file (.css)",
	"font_family":    "Application font family",
	"font_size":      "Application font size",

	"search_on_type":            "Enable/disable search-on-type",
	"search_on_type_min_length": "Minimum query length for search-on-type",

	"header_template": "HTML template for header (dict name + entry terms)",

	"history_disable":   "Disable history",
	"history_auto_save": "Auto-save history on every new record",
	"history_max_size":  "Maximum size for history",

	"most_frequent_disable":   "Disable keeping Most Frequent queries",
	"most_frequent_auto_save": "Auto-save Most Frequent queries",
	"most_frequent_max_size":  "Maximum size for Most Frequent queries",

	"favorites_auto_save": "Auto-save Favorites on every new record",

	"max_results_total": "Maximum number of search results",

	"audio": "Enable audio in article",

	"embed_external_stylesheet": "Embed external stylesheet/css in article",

	"color_mapping": "Mapping for colors used in article",

	"popup_style_str": "Stylesheet (text) for 'Loading' popup",

	"article_zoom_factor": "Zoom factor for article with mouse wheel or keyboard",

	"article_arrow_keys": "Use arrow keys to scroll through article (when focused)",

	"reduce_minimum_window_width": "Use smaller buttons to reduce minimum width of window",

	"local_server_ports": "Ports for local server. Server runs on first port; Client tries all",

	"local_client_timeout": "Timeout for local web client, default is 100ms",
}

func getTomlTag(s string) string {
	i := strings.Index(s, `toml:"`)
	if i < 0 {
		return ""
	}
	start := i + 6
	len := strings.Index(s[start:], `"`)
	if len < 0 {
		return ""
	}
	return s[start : start+len]
}

func printCommentTemplate() {
	conf := config.Default()
	typ := reflect.TypeOf(conf).Elem()
	for i := 0; i < typ.NumField(); i++ {
		fieldType := typ.Field(i)
		tomlTag := getTomlTag(string(fieldType.Tag))
		fmt.Printf("%#v: \"\",\n", tomlTag)
	}
}

func printAll() {
	conf := config.Default()
	typ := reflect.TypeOf(conf).Elem()
	val := reflect.ValueOf(conf).Elem()
	for i := 0; i < typ.NumField(); i++ {
		fieldType := typ.Field(i)
		name := fieldType.Name
		key := getTomlTag(string(fieldType.Tag))
		fieldVal := val.Field(i)
		fieldValIn := fieldVal.Interface()
		comment := commentMap[key]
		fmt.Printf(
			"name=%v, toml=%v, default=%#v, comment=%#v\n\n",
			name,
			key,
			fieldValIn,
			comment,
		)
	}
}

func printMarkdown() {
	conf := config.Default()
	typ := reflect.TypeOf(conf).Elem()
	val := reflect.ValueOf(conf).Elem()
	for i := 0; i < typ.NumField(); i++ {
		fieldType := typ.Field(i)
		key := getTomlTag(string(fieldType.Tag))
		fieldVal := val.Field(i)
		fieldValIn := fieldVal.Interface()
		comment := commentMap[key]

		keyCode := codeValue(key)
		fmt.Println(keyCode)
		fmt.Println(strings.Repeat("-", len(keyCode)))
		fmt.Println(comment + "\n")
		fmt.Println("Default value: " + jsonCodeValue(fieldValIn) + "\n")
	}
}

func main() {
	cmd := ""
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}
	switch cmd {
	case "comment-template":
		printCommentTemplate()
		return
	case "debug":
		printAll()
		return
		// case "gen-table":
		// 	printMarkdownTable()
		// 	return
		// case "gen-list":
		// 	printMarkdownList()
		// 	return
	}
	printMarkdown()
}
