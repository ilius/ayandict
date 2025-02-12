module github.com/ilius/ayandict/v2

go 1.24.0

toolchain go1.24.6

require (
	codeberg.org/ilius/go-dict-commons v0.7.0
	codeberg.org/ilius/go-dict-sql v0.6.0
	codeberg.org/ilius/go-stardict/v2 v2.7.1
	github.com/BurntSushi/toml v1.5.0
	github.com/ilius/is/v2 v2.3.2
	github.com/ilius/qt v0.0.0-20230422004322-c855bcf0151b
)

require (
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gopherjs/gopherjs v1.17.2 // indirect
	github.com/ilius/glob v0.0.0-20250212111036-4c41f838a304 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/ncruces/go-strftime v0.1.9 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	golang.org/x/exp v0.0.0-20251002181428-27f1f14c8bb9 // indirect
	golang.org/x/sys v0.36.0 // indirect
	modernc.org/libc v1.66.10 // indirect
	modernc.org/mathutil v1.7.1 // indirect
	modernc.org/memory v1.11.0 // indirect
	modernc.org/sqlite v1.39.0 // indirect
)

// replace codeberg.org/ilius/go-stardict/v2 => ../go-stardict
// replace codeberg.org/ilius/go-dict-sql => ../go-dict-sql
// replace codeberg.org/ilius/go-dict-commons => ../go-dict-commons
