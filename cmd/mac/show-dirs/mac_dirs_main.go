package main

import (
	"fmt"

	"github.com/ilius/ayandict/v3/pkg/config"
)

func main() {
	fmt.Println("Config:", config.GetConfigDir())
	fmt.Println("Cache:", config.GetCacheDir())
	fmt.Println("State:", config.GetStateDir())
}
