package web

import "embed"

//go:embed web/*.html web/*.css web/*.png web/brython@*/*.js
var FS embed.FS
