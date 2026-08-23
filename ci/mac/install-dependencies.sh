#!/usr/bin/env bash

set -euo pipefail

brew install \
  brotli \
  gcc \
  cmake \
  ninja \
  openssl@3 \
  harfbuzz \
  libpng \
  pcre2 \
  zlib \
  pkg-config \
  sevenzip
