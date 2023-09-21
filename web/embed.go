package web

import "embed"

//go:embed */*.html */*.css */*.png */*/*.js
var FS embed.FS
