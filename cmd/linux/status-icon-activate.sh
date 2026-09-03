#!/usr/bin/env bash
set -o errexit -o pipefail -o nounset

if command -v socat >/dev/null 2>&1 ; then
    printf 'statusicon:activate' | socat - "UNIX-CONNECT:/tmp/ayandict-$UID"
else
    printf 'statusicon:activate' | nc -U "/tmp/ayandict-$UID" -q 2
fi
