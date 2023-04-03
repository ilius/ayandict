package main

import (
	"log"
	"os"

	"github.com/ilius/ayandict/pkg/application"
)

func main() {
	log.SetOutput(os.Stdout)
	application.Run()
}
