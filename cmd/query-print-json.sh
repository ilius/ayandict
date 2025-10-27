#!/bin/bash
set -o errexit -o pipefail -o nounset

# This queries given arguments via socket API and prints the results in json format
echo "query:fuzzy:$@" | socat -t 5 - UNIX-CONNECT:/tmp/ayandict-$UID,crnl
echo