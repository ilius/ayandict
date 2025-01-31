#!/bin/bash
set -e
set -x

if [ -z "$GOARCH" ] ; then
	echo "GOARCH is not set" 2>&1
	exit 1
fi

export CGO_ENABLED=1
export GOOS=darwin

FLAGS=(-ldflags '-s -w' -trimpath)
VERSION=$(go run pkg/version/version.go)

OUT=ayandict-$VERSION-mac-$GOARCH
go build -o $OUT "${FLAGS[@]}" "$@"
bzip2 -f $OUT
