#!/usr/bin/env python3

from __future__ import annotations

import os
import shlex
import shutil
import subprocess
import tempfile
from pathlib import Path

from config import MAC_CI_DIR, load_config


def run(
	command: list[str],
	*,
	env: dict[str, str],
	cwd: Path | None = None,
) -> None:
	subprocess.run(command, env=env, cwd=cwd, check=True)


def generated_objects(link_parts: list[str], build_dir: Path) -> list[Path]:
	objects: list[Path] = []
	for argument in link_parts:
		if not argument.endswith(".o") or argument.endswith("main.cpp.o"):
			continue
		path = Path(argument)
		objects.append(path if path.is_absolute() else build_dir / path)
	if not objects:
		raise RuntimeError("Qt's CMake link command contained no generated objects")
	return objects


def linker_flags_after_main(link_parts: list[str]) -> list[str]:
	try:
		main_index = next(
			index
			for index, argument in enumerate(link_parts)
			if argument.endswith("CMakeFiles/qt_link_probe.dir/main.cpp.o")
		)
	except StopIteration as error:
		raise RuntimeError("Could not extract the Qt CMake link interface") from error

	flags: list[str] = []
	skip_output = False
	for argument in link_parts[main_index + 1 :]:
		if argument in {"&&", ":"}:
			continue
		if argument == "-o":
			skip_output = True
			continue
		if skip_output:
			skip_output = False
			continue
		if not argument.endswith(".o"):
			flags.append(argument)
	return flags


def main() -> None:
	config = load_config()
	repo_root = MAC_CI_DIR.parent.parent
	environment = os.environ.copy()
	environment.update(
		{
			"CGO_ENABLED": "1",
			"PKG_CONFIG_PATH": str(config.qt_prefix / "lib/pkgconfig"),
			"GOCACHE": str(config.go_cache_dir),
			"CGO_CXXFLAGS": "-std=c++17",
		}
	)
	library_path = str(config.qt_prefix / "lib")
	if existing_library_path := environment.get("LIBRARY_PATH"):
		library_path = f"{library_path}:{existing_library_path}"
	environment["LIBRARY_PATH"] = library_path
	environment["CGO_CPPFLAGS"] = " ".join(
		(
			f"-I{config.qt_prefix}/include",
			f"-F{config.qt_prefix}/lib",
			"-iframework",
			f"{config.qt_prefix}/lib",
			"-isystem",
			f"{config.qt_prefix}/lib/QtCore.framework/Headers",
			"-isystem",
			f"{config.qt_prefix}/lib/QtGui.framework/Headers",
			"-isystem",
			f"{config.qt_prefix}/lib/QtWidgets.framework/Headers",
			"-isystem",
			f"{config.qt_prefix}/lib/QtNetwork.framework/Headers",
			"-isystem",
			f"{config.qt_prefix}/lib/QtMultimedia.framework/Headers",
			"-isystem",
			f"{config.qt_prefix}/lib/QtMultimediaWidgets.framework/Headers",
		)
	)
	config.go_cache_dir.mkdir(parents=True, exist_ok=True)

	with tempfile.TemporaryDirectory(
		prefix="ayandict-qt-link-probe.", dir="/tmp"
	) as probe_directory:
		probe = Path(probe_directory)
		build_dir = probe / "build"
		run(
			[
				"cmake",
				"-S",
				str(MAC_CI_DIR / "qt-link-probe"),
				"-B",
				str(build_dir),
				"-GNinja",
				"-DCMAKE_BUILD_TYPE=Release",
				f"-DCMAKE_PREFIX_PATH={config.qt_prefix}",
			],
			env=environment,
		)
		run(["cmake", "--build", str(build_dir), "--verbose"], env=environment)

		commands = subprocess.run(
			["ninja", "-C", str(build_dir), "-t", "commands", "qt_link_probe"],
			env=environment,
			check=True,
			stdout=subprocess.PIPE,
			text=True,
		).stdout.splitlines()
		if not commands:
			raise RuntimeError("Ninja returned no Qt link command")
		link_parts = shlex.split(commands[-1])

		objects = generated_objects(link_parts, build_dir)
		generated_archive = probe / "libqt-generated.a"
		run(
			[
				"/usr/bin/libtool",
				"-static",
				"-o",
				str(generated_archive),
				*(str(path) for path in objects),
			],
			env=environment,
		)

		real_pkg_config = shutil.which("pkg-config")
		if not real_pkg_config:
			raise RuntimeError("pkg-config is not installed")
		environment["REAL_PKG_CONFIG"] = real_pkg_config
		environment["PKG_CONFIG"] = str(MAC_CI_DIR / "pkg-config-for-cgo.py")
		environment["CGO_LDFLAGS"] = " ".join(
			(
				f"-L{config.qt_prefix}/lib",
				f"-Wl,-force_load,{generated_archive}",
				*linker_flags_after_main(link_parts),
			)
		)
		run(
			["go", "build", "-ldflags", "-s -w", "-trimpath"],
			env=environment,
			cwd=repo_root,
		)


if __name__ == "__main__":
	main()
