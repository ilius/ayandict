#!/usr/bin/env bash

set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=ci/mac/config.sh
source "$script_dir/config.sh"

mkdir -p "$QT_ROOT"
curl --fail --location --retry 3 \
  --output "$QT_ROOT/$QT_SOURCE_ARCHIVE" \
  "$QT_SOURCE_URL"
tar -xf "$QT_ROOT/$QT_SOURCE_ARCHIVE" -C "$QT_ROOT"
cd "$QT_SOURCE_DIR"

# Qt normally generates pkg-config files only for shared builds, while MIQT
# requires them through its #cgo pkg-config directives.
sed -i '' \
  '/^[[:space:]]*if(NOT BUILD_SHARED_LIBS)$/,/^[[:space:]]*endif()$/d' \
  qtbase/cmake/QtPkgConfigHelpers.cmake

./configure \
  -prefix "$QT_PREFIX" \
  -release \
  -static \
  -static-runtime \
  -opensource \
  -confirm-license \
  -nomake examples \
  -nomake tests \
  -no-opengl \
  -no-dbus \
  -no-gtk \
  -no-feature-zstd \
  -skip qtwebengine -skip qtquick3d -skip qtgraphs \
  -skip qtquick3dphysics -skip qtcanvaspainter -skip qtdeclarative \
  -skip qtquicktimeline -skip qtdoc -skip qtlocation -skip qtlottie \
  -skip qtmqtt -skip qtopcua -skip qtquickeffectmaker \
  -skip qtvirtualkeyboard -skip qtwebview \
  -- \
  -DCMAKE_POLICY_DEFAULT_CMP0177=NEW \
  -DFEATURE_pkg_config=ON

cmake --build . --parallel "$(sysctl -n hw.logicalcpu)"
cmake --install .
