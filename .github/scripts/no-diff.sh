#!/usr/bin/env bash
set -o errexit -o nounset -o pipefail

CHANGES=$(git diff --name-only HEAD --)
if [ -n "$CHANGES" ]; then
	echo "There are changes after running gen.sh:"
	echo "$CHANGES"
	git diff
	exit 1
fi
