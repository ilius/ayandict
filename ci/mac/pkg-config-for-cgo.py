#!/usr/bin/env python3

from __future__ import annotations

import os
import re
import subprocess
import sys


def main() -> None:
    real_pkg_config = os.environ.get("REAL_PKG_CONFIG")
    if not real_pkg_config:
        raise RuntimeError("REAL_PKG_CONFIG must point to the real pkg-config")

    result = subprocess.run(
        [real_pkg_config, *sys.argv[1:]],
        check=False,
        stdout=subprocess.PIPE,
    )
    filtered = re.sub(rb"(^|\s)-lQt6\S+", b"", result.stdout)
    sys.stdout.buffer.write(filtered)
    raise SystemExit(result.returncode)


if __name__ == "__main__":
    main()
