#!/usr/bin/env bash

set -euo pipefail

: "${REAL_PKG_CONFIG:?REAL_PKG_CONFIG must point to the real pkg-config}"

"$REAL_PKG_CONFIG" "$@" |
  sed -E 's#(^|[[:space:]])-lQt6[^[:space:]]+##g'
