module github.com/ilius/ayandict

go 1.19

require (
	github.com/BurntSushi/toml v1.2.1
	github.com/ilius/go-stardict v0.3.1-0.20230322021424-a3ed5e190dcc
	github.com/therecipe/qt v0.0.0-20200904063919-c0c124a5770d
)

require github.com/gopherjs/gopherjs v0.0.0-20190411002643-bd77b112433e // indirect

replace github.com/ilius/go-stardict => ../go-stardict
