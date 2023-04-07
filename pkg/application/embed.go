package application

import "embed"

//go:embed res/*.png res/*.svg
var res embed.FS
