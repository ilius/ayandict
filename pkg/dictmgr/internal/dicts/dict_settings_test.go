package dicts

import (
	"testing"

	"github.com/ilius/is/v2"
)

func TestDictSettingsFuzzy(t *testing.T) {
	is := is.New(t)
	{
		ds := &DictionarySettings{}
		is.True(ds.Fuzzy())
	}
	{
		ds := &DictionarySettings{}
		ds.SetFuzzy(true)
		is.True(ds.Fuzzy())
	}
	{
		ds := &DictionarySettings{}
		ds.SetFuzzy(false)
		is.False(ds.Fuzzy())
	}
}

func TestDictSettingsStartWith(t *testing.T) {
	is := is.New(t)
	{
		ds := &DictionarySettings{}
		is.True(ds.StartWith())
	}
	{
		ds := &DictionarySettings{}
		ds.SetStartWith(true)
		is.True(ds.StartWith())
	}
	{
		ds := &DictionarySettings{}
		ds.SetStartWith(false)
		is.False(ds.StartWith())
	}
}

func TestDictSettingsRegex(t *testing.T) {
	is := is.New(t)
	{
		ds := &DictionarySettings{}
		is.True(ds.Regex())
	}
	{
		ds := &DictionarySettings{}
		ds.SetRegex(true)
		is.True(ds.Regex())
	}
	{
		ds := &DictionarySettings{}
		ds.SetRegex(false)
		is.False(ds.Regex())
	}
}

func TestDictSettingsGlob(t *testing.T) {
	is := is.New(t)
	{
		ds := &DictionarySettings{}
		is.True(ds.Glob())
	}
	{
		ds := &DictionarySettings{}
		ds.SetGlob(true)
		is.True(ds.Glob())
	}
	{
		ds := &DictionarySettings{}
		ds.SetGlob(false)
		is.False(ds.Glob())
	}
}

func TestDictSettingsFlagsMixed(t *testing.T) {
	is := is.New(t)
	{
		ds := &DictionarySettings{}
		is.True(ds.Fuzzy())
		is.True(ds.StartWith())
		is.True(ds.Regex())
		is.True(ds.Glob())
		ds.SetFuzzy(false)
		is.False(ds.Fuzzy())
		is.True(ds.StartWith())
		is.True(ds.Regex())
		is.True(ds.Glob())
		ds.SetGlob(false)
		is.False(ds.Fuzzy())
		is.True(ds.StartWith())
		is.True(ds.Regex())
		is.False(ds.Glob())
		ds.SetRegex(false)
		is.False(ds.Fuzzy())
		is.True(ds.StartWith())
		is.False(ds.Regex())
		is.False(ds.Glob())
		ds.SetStartWith(false)
		is.False(ds.Fuzzy())
		is.False(ds.StartWith())
		is.False(ds.Regex())
		is.False(ds.Glob())
	}
	{
		ds := &DictionarySettings{}
		is.True(ds.Fuzzy())
		is.True(ds.StartWith())
		is.True(ds.Regex())
		is.True(ds.Glob())
		ds.SetStartWith(false)
		is.True(ds.Fuzzy())
		is.False(ds.StartWith())
		is.True(ds.Regex())
		is.True(ds.Glob())
	}
}
