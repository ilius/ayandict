package main

import (
	"log"
	"os"

	"github.com/ilius/ayandict/pkg/application"
	"github.com/ilius/ayandict/pkg/qerr"
	"github.com/ilius/ayandict/pkg/stardict"
)

func main() {
	log.SetOutput(os.Stdout)
	stardict.ErrorHandler = func(err error) {
		qerr.Error(err)
	}
	application.Run()
}
