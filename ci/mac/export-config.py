#!/usr/bin/env python3

from __future__ import annotations

import os
import sys
from pathlib import Path

from config import load_config


def main() -> None:
    config = load_config()
    output = "".join(
        (
            f"qt-version={config.qt_version}\n",
            f"qt-prefix={config.qt_prefix}\n",
            f"qt-cache-key={config.qt_cache_key}\n",
            f"go-version={config.go_version}\n",
            f"go-cache-root={config.go_cache_root}\n",
        )
    )

    github_output = os.environ.get("GITHUB_OUTPUT")
    if github_output:
        with Path(github_output).open("a", encoding="utf-8") as output_file:
            output_file.write(output)
    else:
        sys.stdout.write(output)


if __name__ == "__main__":
    main()
