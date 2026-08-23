#!/usr/bin/env bash

set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=ci/mac/config.sh
source "$script_dir/config.sh"

output_file="${GITHUB_OUTPUT:-/dev/stdout}"
{
  printf 'qt-version=%s\n' "$QT_VERSION"
  printf 'qt-prefix=%s\n' "$QT_PREFIX"
  printf 'qt-cache-key=%s\n' "$QT_CACHE_KEY"
  printf 'go-version=%s\n' "$GO_VERSION"
  printf 'go-cache-root=%s\n' "$GO_CACHE_ROOT"
} >> "$output_file"
