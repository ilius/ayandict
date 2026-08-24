#!/usr/bin/env bash
set -o errexit -o pipefail -o nounset

DIR=$(dirname -- "$0")
cd "$DIR"
./scan-popup-basic "$(xclip -o)"
