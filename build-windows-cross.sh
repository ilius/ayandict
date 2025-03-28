#!/bin/bash
set -e
set -x

FLAGS=(-ldflags '-s -w' -trimpath)
VERSION=$(go run pkg/version/version.go)

TARGET=x86_64-windows
export CC="zig cc -target $TARGET"
export CXX="zig c++ -target $TARGET"

export CGO_ENABLED=1
export GOOS=windows
export GOARCH=amd64
OUT=ayandict-$VERSION-windows-amd64.exe
time go build -o $OUT "${FLAGS[@]}" "$@"

