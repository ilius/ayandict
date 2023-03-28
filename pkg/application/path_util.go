package application

import (
	"path/filepath"
	"runtime"
)

func PathFromUnix(pathStr string) string {
	if runtime.GOOS != "windows" {
		return pathStr
	}
	if pathStr == "" {
		return ""
	}
	pathStr = filepath.FromSlash(pathStr)
	if pathStr[0] != '\\' {
		return pathStr
	}
	if pathStr[2] != '\\' {
		return pathStr
	}
	// change `\C\Users` to `C:\Users`
	pathStr = pathStr[1:2] + ":\\" + pathStr[3:]
	return pathStr
}
