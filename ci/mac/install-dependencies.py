#!/usr/bin/env python3

from __future__ import annotations

import subprocess

DEPENDENCIES = (
	"cmake",
	"ninja",
	"pkg-config",
	"sevenzip",
)


def main() -> None:
	subprocess.run(["brew", "install", *DEPENDENCIES], check=True)


if __name__ == "__main__":
	main()
