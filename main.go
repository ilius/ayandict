package main

import (
	"log"
	"os"

	"github.com/ilius/ayandict/v2/pkg/application"
)

func main() {
	log.SetOutput(os.Stdout)

	application.Run()
}
