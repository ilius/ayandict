#!/bin/bash
set -o errexit -o pipefail -o nounset

echo statusicon:activate | socat - UNIX-CONNECT:/tmp/ayandict-$UID
