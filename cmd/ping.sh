#!/usr/bin/env bash
set -o errexit -o pipefail -o nounset

echo ping | socat -t 5 - UNIX-CONNECT:/tmp/ayandict-$UID
echo