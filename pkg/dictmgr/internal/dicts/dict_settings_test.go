package dicts

import (
	"testing"

	"github.com/ilius/is/v2"
)

func TestDictSettingsFuzzy(t *testing.T) {
	is := is.New(t)
	{
		ds := &DictSettings{}
		is.True(ds.Fuzzy())
	}
	{
		ds := &DictSettings{}
		ds.SetFuzzy(true)
		is.True(ds.Fuzzy())
	}
	{
		ds := &DictSettings{}
		ds.SetFuzzy(false)
		is.False(ds.Fuzzy())
	}
}

func TestDictSettingsStartWith(t *testing.T) {
	is := is.New(t)
	{
		ds := &DictSettings{}
		is.True(ds.StartWith())
	}
	{
		ds := &DictSettings{}
		ds.SetStartWith(true)
		is.True(ds.StartWith())
	}
	{
		ds := &DictSettings{}
		ds.SetStartWith(false)
		is.False(ds.StartWith())
	}
}

func TestDictSettingsRegex(t *testing.T) {
	is := is.New(t)
	{
		ds := &DictSettings{}
		is.True(ds.Regex())
	}
	{
		ds := &DictSettings{}
		ds.SetRegex(true)
		is.True(ds.Regex())
	}
	{
		ds := &DictSettings{}
		ds.SetRegex(false)
		is.False(ds.Regex())
	}
}

func TestDictSettingsGlob(t *testing.T) {
	is := is.New(t)
	{
		ds := &DictSettings{}
		is.True(ds.Glob())
	}
	{
		ds := &DictSettings{}
		ds.SetGlob(true)
		is.True(ds.Glob())
	}
	{
		ds := &DictSettings{}
		ds.SetGlob(false)
		is.False(ds.Glob())
	}
}

func TestDictSettingsFlagsMixed(t *testing.T) {
	is := is.New(t)
	{
		ds := &DictSettings{}
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
		ds := &DictSettings{}
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
