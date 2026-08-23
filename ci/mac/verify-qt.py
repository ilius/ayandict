#!/usr/bin/env python3

from __future__ import annotations

import os
import subprocess

from config import load_config

QT_PACKAGES = ("Qt6Widgets", "Qt6Network", "Qt6MultimediaWidgets")


def main() -> None:
	config = load_config()
	environment = os.environ.copy()
	environment["PKG_CONFIG_PATH"] = str(config.qt_prefix / "lib/pkgconfig")

	subprocess.run(["ls", "-l", str(config.qt_prefix / "lib")], check=True)
	subprocess.run(
		[
			"file",
			str(config.qt_prefix / "lib/QtCore.framework/Versions/A/QtCore"),
		],
		check=True,
	)
	commands = (
		["pkg-config", "--modversion", *QT_PACKAGES],
		["pkg-config", "--cflags", *QT_PACKAGES],
		["pkg-config", "--libs", "--static", *QT_PACKAGES],
	)
	for command in commands:
		subprocess.run(command, env=environment, check=True)


if __name__ == "__main__":
	main()
