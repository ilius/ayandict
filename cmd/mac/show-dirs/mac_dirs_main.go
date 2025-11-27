package main

import (
	"fmt"

	"github.com/ilius/ayandict/v3/pkg/config"
)

func main() {
	fmt.Println("Config:", config.Paths.ConfigDir())
	fmt.Println("Cache:", config.Paths.CacheDir())
	fmt.Println("State:", config.Paths.StateDir())
}
