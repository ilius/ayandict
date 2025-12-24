#!/usr/bin/env bash
set -e
set -x

export CGO_ENABLED=1

FLAGS=(-ldflags '-s -w' -trimpath)
VERSION=$(go run pkg/version/version.go)


function run_zip() {
	IN_PATH=$1
	ZIP_PATH="${IN_PATH%.*}.zip"
	if [ -f C:\\Windows\\System32\\tar.exe ] ; then
		C:\\Windows\\System32\\tar.exe -a -c -f $ZIP_PATH $IN_PATH
		rm $IN_PATH
		return
	fi
	if which zip ; then
		zip $ZIP_PATH $IN_PATH
		return
	fi
}


OUT=ayandict-$VERSION-windows-$(go env GOARCH).exe
go build -o $OUT "${FLAGS[@]}" "$@"
run_zip $OUT

