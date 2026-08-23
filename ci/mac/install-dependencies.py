#!/usr/bin/env python3

from __future__ import annotations

import subprocess

DEPENDENCIES = (
    "brotli",
    "gcc",
    "cmake",
    "ninja",
    "openssl@3",
    "harfbuzz",
    "libpng",
    "pcre2",
    "zlib",
    "pkg-config",
    "sevenzip",
)


def main() -> None:
    subprocess.run(["brew", "install", *DEPENDENCIES], check=True)


if __name__ == "__main__":
    main()
