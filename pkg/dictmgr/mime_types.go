package dictmgr

var extByMimeTypeMap = map[string]string{
	"image/png":     "png",
	"image/jpeg":    "jpg",
	"image/gif":     "gif",
	"image/svg+xml": "svg",
	"image/webp":    "webp",
	"image/tiff":    "tiff",
	"image/bmp":     "bmp",
	"image/x-icon":  "ico",

	"text/css":        "css",
	"text/plain":      "ini",
	"text/javascript": "js",

	"application/javascript":        "js",
	"application/json":              "json",
	"application/pdf":               "pdf",
	"application/font-woff":         "woff",
	"application/font-woff2":        "woff2",
	"application/x-font-ttf":        "ttf",
	"application/x-font-opentype":   "otf",
	"application/vnd.ms-fontobject": "eot",

	"audio/mpeg":    "mp3",
	"audio/ogg":     "ogg",
	"audio/x-speex": "spx",
	"audio/wav":     "wav",
	"video/mp4":     "mp4",

	// "application/octet-stream+xapian": "",
	// "application/octet-stream": "",
	// "application/x-chrome-extension",
	// "application/warc-headers",
}
