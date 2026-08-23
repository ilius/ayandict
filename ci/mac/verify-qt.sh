#!/usr/bin/env bash

set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=ci/mac/config.sh
source "$script_dir/config.sh"

export PKG_CONFIG_PATH="$QT_PREFIX/lib/pkgconfig"

ls -l "$QT_PREFIX/lib"
file "$QT_PREFIX/lib/QtCore.framework/Versions/A/QtCore"
pkg-config --modversion Qt6Widgets Qt6Network Qt6MultimediaWidgets
pkg-config --cflags Qt6Widgets Qt6Network Qt6MultimediaWidgets
pkg-config --libs --static Qt6Widgets Qt6Network Qt6MultimediaWidgets
