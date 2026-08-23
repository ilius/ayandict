#!/usr/bin/env bash

set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=ci/mac/config.sh
source "$script_dir/config.sh"

repo_root="$(cd "$script_dir/../.." && pwd)"
cd "$repo_root"

export CGO_ENABLED=1
export PKG_CONFIG_PATH="$QT_PREFIX/lib/pkgconfig"
export GOCACHE="$GO_CACHE_DIR"
export LIBRARY_PATH="$QT_PREFIX/lib${LIBRARY_PATH:+:$LIBRARY_PATH}"
export CGO_CPPFLAGS="\
  -I$QT_PREFIX/include \
  -F$QT_PREFIX/lib \
  -iframework $QT_PREFIX/lib \
  -isystem $QT_PREFIX/lib/QtCore.framework/Headers \
  -isystem $QT_PREFIX/lib/QtGui.framework/Headers \
  -isystem $QT_PREFIX/lib/QtWidgets.framework/Headers \
  -isystem $QT_PREFIX/lib/QtNetwork.framework/Headers \
  -isystem $QT_PREFIX/lib/QtMultimedia.framework/Headers \
  -isystem $QT_PREFIX/lib/QtMultimediaWidgets.framework/Headers"
export CGO_CXXFLAGS="-std=c++17"

mkdir -p "$GOCACHE"
qt_probe="$(mktemp -d /tmp/ayandict-qt-link-probe.XXXXXX)"
cleanup() {
  rm -rf "$qt_probe"
}
trap cleanup EXIT

cmake -S "$MAC_CI_DIR/qt-link-probe" -B "$qt_probe/build" -GNinja \
  -DCMAKE_BUILD_TYPE=Release \
  -DCMAKE_PREFIX_PATH="$QT_PREFIX"
cmake --build "$qt_probe/build" --verbose

qt_link_line="$(ninja -C "$qt_probe/build" -t commands qt_link_probe | tail -n 1)"
# Match the object by suffix: Ninja may print a relative or absolute path.
qt_link_args="${qt_link_line#*CMakeFiles/qt_link_probe.dir/main.cpp.o }"
if [[ "$qt_link_args" == "$qt_link_line" ]]; then
  echo 'Could not extract the Qt CMake link interface' >&2
  exit 1
fi

read -r -a qt_link_parts <<< "$qt_link_line"
qt_generated_objects=()
for arg in "${qt_link_parts[@]}"; do
  if [[ "$arg" == *.o && "$arg" != *main.cpp.o ]]; then
    if [[ "$arg" != /* ]]; then
      arg="$qt_probe/build/$arg"
    fi
    qt_generated_objects+=("$arg")
  fi
done

qt_generated_archive="$qt_probe/libqt-generated.a"
/usr/bin/libtool -static \
  -o "$qt_generated_archive" \
  "${qt_generated_objects[@]}"

read -r -a qt_link_arg_parts <<< "$qt_link_args"
qt_cgo_ldflags=""
skip_link_output=false
for arg in "${qt_link_arg_parts[@]}"; do
  if [[ "$arg" == "&&" || "$arg" == ":" ]]; then
    continue
  fi
  if [[ "$arg" == "-o" ]]; then
    skip_link_output=true
    continue
  fi
  if [[ "$skip_link_output" == true ]]; then
    skip_link_output=false
    continue
  fi
  [[ "$arg" == *.o ]] || qt_cgo_ldflags+=" $arg"
done

# MIQT invokes pkg-config itself. Its redundant -lQt6* flags do not understand
# macOS static frameworks; CMake already supplied the exact framework binaries.
REAL_PKG_CONFIG="$(command -v pkg-config)"
export REAL_PKG_CONFIG
export PKG_CONFIG="$MAC_CI_DIR/pkg-config-for-cgo.sh"
export CGO_LDFLAGS="\
  -L$QT_PREFIX/lib \
  -Wl,-force_load,$qt_generated_archive \
  $qt_cgo_ldflags"

go build -ldflags '-s -w' -trimpath
