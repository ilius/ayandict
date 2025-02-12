module github.com/ilius/ayandict/v2

go 1.22

require (
	github.com/BurntSushi/toml v1.3.2
	github.com/ilius/go-dict-commons v0.4.2
	github.com/ilius/go-stardict/v2 v2.3.1
	github.com/ilius/is/v2 v2.3.2
	github.com/ilius/qt v0.0.0-20230422004322-c855bcf0151b
	golang.org/x/sys v0.17.0
)

require (
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/gopherjs/gopherjs v1.17.2 // indirect
)

// replace github.com/ilius/go-stardict/v2 => ../go-stardict
// replace github.com/ilius/go-dict-sql => ../go-dict-sql
// replace github.com/ilius/go-dict-commons => ../go-dict-commons
