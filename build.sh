#!/bin/bash
set -e
set -x

if [ -z "$GOARCH" ] ; then
	echo "GOARCH is not set" 2>&1
	exit 1
fi
if [ -z "$GOOS" ] ; then
	echo "GOOS is not set" 2>&1
	exit 1
fi

export CGO_ENABLED=1
FLAGS=(-ldflags '-s -w' -trimpath)
VERSION=$(go run pkg/version/version.go)

OUT=ayandict-$VERSION-$GOOS-$GOARCH
time go build -o $OUT "${FLAGS[@]}" "$@"
ls -lh $OUT

bzip2 -f "$OUT"

ls -lh $OUT.bz2


