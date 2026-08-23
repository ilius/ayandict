#!/usr/bin/env python3

from __future__ import annotations

import re

from config import load_config


def main() -> None:
	config = load_config()
	pkg_config_dir = config.qt_prefix / "lib/pkgconfig"
	pc_files = sorted(pkg_config_dir.glob("Qt6*.pc"))
	if not pc_files:
		raise RuntimeError(f"No Qt pkg-config files found in {pkg_config_dir}")

	for pc_file in pc_files:
		package_name = pc_file.stem
		framework_name = f"Qt{package_name.removeprefix('Qt6')}"
		framework_binary = (
			config.qt_prefix
			/ "lib"
			/ f"{framework_name}.framework"
			/ "Versions/A"
			/ framework_name
		)
		if not framework_binary.is_file():
			continue

		contents = pc_file.read_text(encoding="utf-8")
		pattern = re.compile(rf"-l{re.escape(package_name)}(?=\s|$)")
		updated = pattern.sub(str(framework_binary), contents)
		if updated != contents:
			pc_file.write_text(updated, encoding="utf-8")


if __name__ == "__main__":
	main()
