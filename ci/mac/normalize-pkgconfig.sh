#!/usr/bin/env bash

set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=ci/mac/config.sh
source "$script_dir/config.sh"

shopt -s nullglob
pc_files=("$QT_PREFIX"/lib/pkgconfig/Qt6*.pc)
if (( ${#pc_files[@]} == 0 )); then
  echo "No Qt pkg-config files found in $QT_PREFIX/lib/pkgconfig" >&2
  exit 1
fi

for pc_file in "${pc_files[@]}"; do
  package_name="$(basename "$pc_file" .pc)"
  framework_name="Qt${package_name#Qt6}"
  framework_binary="$QT_PREFIX/lib/$framework_name.framework/Versions/A/$framework_name"
  if [[ -f "$framework_binary" ]]; then
    sed -E -i '' \
      "s#-l${package_name}([[:space:]]|$)#${framework_binary}\\1#g" \
      "$pc_file"
  fi
done
