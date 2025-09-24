package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/ilius/ayandict/v2/pkg/config"
)

func getFieldTag(s string, tag string) string {
	i := strings.Index(s, tag+`:"`)
	if i < 0 {
		return ""
	}
	start := i + len(tag) + 2
	len := strings.Index(s[start:], `"`)
	if len < 0 {
		return ""
	}
	return s[start : start+len]
}

func getTomlTag(s string) string {
	return getFieldTag(s, "toml")
}

func getDocTag(s string) string {
	return strings.ReplaceAll(getFieldTag(s, "doc"), "‘", "`")
}

func printCommentTemplate() {
	conf := config.Default()
	typ := reflect.TypeOf(conf).Elem()
	for i := range typ.NumField() {
		fieldType := typ.Field(i)
		tomlTag := getTomlTag(string(fieldType.Tag))
		fmt.Printf("%#v: \"\",\n", tomlTag)
	}
}

func printStruct(typ reflect.Type, val reflect.Value, keyPrefix string) {
	for i := range typ.NumField() {
		field := typ.Field(i)
		name := field.Name
		key := keyPrefix + getTomlTag(string(field.Tag))
		fieldVal := val.Field(i)
		fieldValIn := fieldVal.Interface()
		comment := getDocTag(string(field.Tag))
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
		fieldType := field.Type
		if fieldType.Kind() == reflect.Struct {
			printStruct(fieldType, fieldVal, key+".")
		}
	}
}

func printAll() {
	conf := config.Default()
	typ := reflect.TypeOf(conf).Elem()
	val := reflect.ValueOf(conf).Elem()
	printStruct(typ, val, "")
}

type ConfigStructSpec struct {
	Type    reflect.Type
	Value   reflect.Value
	KeyPath []string
}

func printMarkdownStruct(spec ConfigStructSpec) {
	subStructs := []ConfigStructSpec{}
	for i := range spec.Type.NumField() {
		field := spec.Type.Field(i)
		keyPath := append(spec.KeyPath, getTomlTag(string(field.Tag)))

		fieldType := field.Type
		fieldValue := spec.Value.Field(i)

		if fieldType.Kind() == reflect.Struct {
			subStructs = append(subStructs, ConfigStructSpec{
				Type:    fieldType,
				Value:   fieldValue,
				KeyPath: keyPath,
			})
			continue
		}

		comment := getDocTag(string(field.Tag))
		if comment == "" {
			log.Fatalln("No comment for", field)
		}
		comment = strings.ReplaceAll(comment, "`", "``")

		keyCode := codeValue(strings.Join(keyPath, "."))
		fmt.Println(keyCode)
		fmt.Println(strings.Repeat("-", len(keyCode)))
		fmt.Println(comment + "\n")
		fmt.Println("Default value: " + jsonCodeValue(fieldValue.Interface()) + "\n")

	}
	for _, sub := range subStructs {
		printMarkdownStruct(sub)
	}
}

func printMarkdown() {
	conf := config.Default()
	printMarkdownStruct(ConfigStructSpec{
		Type:  reflect.TypeOf(conf).Elem(),
		Value: reflect.ValueOf(conf).Elem(),
	})
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
