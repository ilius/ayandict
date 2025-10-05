package application

import (
	"log"
	"os"
)

var stdLogger = log.New(os.Stderr, "", log.LstdFlags)
