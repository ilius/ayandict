module github.com/ilius/ayandict/v3

go 1.24

require (
	codeberg.org/ilius/go-dict-commons v0.8.0
	codeberg.org/ilius/go-stardict/v2 v2.8.0
	github.com/BurntSushi/toml v1.6.0
	github.com/ilius/is/v2 v2.3.2
	github.com/mappu/miqt v0.13.0
	github.com/xrash/smetrics v0.0.0-20250705151800-55b8f293f342
)

require github.com/ilius/glob v0.0.0-20250212111036-4c41f838a304 // indirect

// replace codeberg.org/ilius/go-stardict/v2 => ../go-stardict
// replace codeberg.org/ilius/go-dict-commons => ../go-dict-commons
// replace github.com/ilius/go-dict-sql => ../go-dict-sql
