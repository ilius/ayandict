#!/usr/bin/env bash
set -o errexit -o pipefail -o nounset

QUERY="$*"
SOCKET="/tmp/ayandict-$UID"
printf 'Opening AyanDict main window with query: %s\n' "$QUERY"
printf 'mainquery:%s' "$QUERY" | socat -t 5 - "UNIX-CONNECT:$SOCKET"
