package ospaths

import (
	"os"
	"testing"
)

func TestPath(t *testing.T) {
	p := &Paths{
		Home:       os.Getenv("HOME"),
		AppName:    "testapp",
		AppNameCap: "TestApp",
		AppDesc:    "Test Application",
	}
	t.Log(p.ConfigDir())
	t.Log(p.CacheDir())
	t.Log(p.StateDir())
}
