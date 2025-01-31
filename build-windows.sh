#!/bin/bash
set -e
set -x

if [ -z "$GOARCH" ] ; then
	echo "GOARCH is not set" 2>&1
	exit 1
fi

export CGO_ENABLED=1
export GOOS=windows

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
		rm $IN_PATH
		return
	fi
}


OUT=ayandict-$VERSION-windows-$GOARCH.exe
go build -o $OUT "${FLAGS[@]}" "$@"
run_zip $OUT

