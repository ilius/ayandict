package main

import (
	"bytes"
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/ilius/ayandict/v2/pkg/config"
)

func main() {
	conf := config.Default()
	buf := bytes.NewBuffer(nil)
	encoder := toml.NewEncoder(buf)
	err := encoder.Encode(conf)
	if err != nil {
		panic(err)
	}
	fmt.Println(buf.String())
}
