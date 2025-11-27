package config

import (
	"os"

	"github.com/ilius/ayandict/v3/pkg/appinfo"
	"github.com/ilius/ayandict/v3/pkg/ospaths"
)

var Paths = ospaths.Paths{
	Home:       os.Getenv("HOME"),
	AppName:    appinfo.APP_NAME,
	AppNameCap: appinfo.APP_DESC,
	AppDesc:    appinfo.APP_DESC,
}
