package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/ilius/ayandict/pkg/config"
)

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
		if comment == "" {
			log.Fatalln("No comment for", key)
		}
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
		if comment == "" {
			log.Fatalln("No comment for", key)
		}

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
