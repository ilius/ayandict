#!/usr/bin/env bash
# shellcheck disable=SC2034 # This file exports configuration to sourcing scripts.

set -euo pipefail

MAC_CI_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
MAC_CONFIG_FILE="$MAC_CI_DIR/config.toml"

toml_string() {
  local key="$1"
  local value

  value="$(awk -v wanted="$key" '
    /^[[:space:]]*#/ { next }
    {
      separator = index($0, "=")
      if (separator == 0) {
        next
      }
      name = substr($0, 1, separator - 1)
      gsub(/^[[:space:]]+|[[:space:]]+$/, "", name)
      if (name != wanted) {
        next
      }
      result = substr($0, separator + 1)
      gsub(/^[[:space:]]+|[[:space:]]+$/, "", result)
      sub(/^"/, "", result)
      sub(/"$/, "", result)
      print result
      exit
    }
  ' "$MAC_CONFIG_FILE")"

  if [[ -z "$value" ]]; then
    echo "Missing '$key' in $MAC_CONFIG_FILE" >&2
    return 1
  fi
  printf '%s' "$value"
}

QT_VERSION="$(toml_string qt_version)"
QT_SOURCE_URL="$(toml_string qt_source_url)"
QT_ROOT="$(toml_string qt_root)"
QT_CACHE_REVISION="$(toml_string qt_cache_revision)"
GO_VERSION="$(toml_string go_version)"
GO_CACHE_ROOT="$(toml_string go_cache_root)"

QT_SOURCE_ARCHIVE="${QT_SOURCE_URL##*/}"
QT_SOURCE_DIR="$QT_ROOT/${QT_SOURCE_ARCHIVE%.tar.xz}"
QT_PREFIX="$QT_ROOT/$QT_VERSION-static"
QT_CACHE_KEY="qt$QT_VERSION-static-macos-$QT_CACHE_REVISION"
GO_CACHE_DIR="$GO_CACHE_ROOT/qt-$QT_VERSION-static"

readonly MAC_CI_DIR MAC_CONFIG_FILE
readonly QT_VERSION QT_SOURCE_URL QT_ROOT QT_CACHE_REVISION
readonly QT_SOURCE_ARCHIVE QT_SOURCE_DIR QT_PREFIX QT_CACHE_KEY
readonly GO_VERSION GO_CACHE_ROOT GO_CACHE_DIR
