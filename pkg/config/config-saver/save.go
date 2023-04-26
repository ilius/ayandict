package main

import "github.com/ilius/ayandict/v2/pkg/config"

func main() {
	err := config.Save(&config.Config{
		FontFamily:   "Shabnam",
		FontSize:     17,
		SearchOnType: false,
	})
	if err != nil {
		panic(err)
	}
}
