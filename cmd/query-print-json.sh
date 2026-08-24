#!/usr/bin/env bash
set -o errexit -o pipefail -o nounset

# This queries given arguments via socket API and prints the results in json format
QUERY="$*"
SOCKET="/tmp/ayandict-$UID"
printf 'query:fuzzy:%s' "$QUERY" | socat -t 5 - "UNIX-CONNECT:$SOCKET,crnl"
echo
