#!/usr/bin/env python3

from __future__ import annotations

import json
import os
import shlex
import shutil
import subprocess
import sys
import tempfile
from pathlib import Path

WINDOWS_CI_DIR = Path(__file__).resolve().parent
REPO_ROOT = WINDOWS_CI_DIR.parent.parent
PROBE_TARGET = "qt_link_probe"
PROBE_SOURCE = "main.cpp"
REQUIRED_PLUGINS = (
	"QWindowsIntegrationPlugin",
	"QModernWindowsStylePlugin",
	"QWindowsMediaPlugin",
	"QGifPlugin",
	"QICOPlugin",
	"QJpegPlugin",
	"QSvgPlugin",
	"QWebpPlugin",
	"QTiffPlugin",
)


def run(
	command: list[str],
	*,
	env: dict[str, str],
	cwd: Path | None = None,
) -> None:
	subprocess.run(command, env=env, cwd=cwd, check=True)


def required_tool(name: str, *, directory: Path | None = None) -> Path:
	if directory is not None:
		for suffix in (".exe", ""):
			candidate = directory / f"{name}{suffix}"
			if candidate.is_file():
				return candidate
	path = shutil.which(name)
	if not path:
		raise RuntimeError(f"Required tool is not installed: {name}")
	return Path(path).resolve()


def create_file_api_query(build_dir: Path) -> None:
	query_dir = build_dir / ".cmake/api/v1/query"
	query_dir.mkdir(parents=True, exist_ok=True)
	(query_dir / "codemodel-v2").touch()


def split_fragment(fragment: str) -> list[str]:
	try:
		# CMake emits native Windows paths in File API link fragments. POSIX
		# shlex treats every backslash as an escape and turns, for example,
		# D:\a\_temp\plugin.obj into D:a_tempplugin.obj. Non-POSIX mode keeps
		# those separators; remove only the surrounding quotes it preserves.
		arguments = shlex.split(fragment, posix=False)
	except ValueError as error:
		raise RuntimeError(f"Could not parse CMake link fragment: {fragment}") from error
	return [
		argument[1:-1]
		if len(argument) >= 2 and argument[0] == argument[-1] == '"'
		else argument
		for argument in arguments
	]


def normalized_path(path: Path) -> str:
	return str(path.resolve()).replace("\\", "/")


def is_object_file(argument: str) -> bool:
	return argument.lower().endswith((".obj", ".o"))


def target_reply(build_dir: Path) -> dict[str, object]:
	replies = list(
		(build_dir / ".cmake/api/v1/reply").glob(f"target-{PROBE_TARGET}-*.json")
	)
	if len(replies) != 1:
		raise RuntimeError(
			f"Expected one CMake File API reply for {PROBE_TARGET}, found {len(replies)}"
		)
	with replies[0].open(encoding="utf-8") as stream:
		return json.load(stream)


def command_fragments(target: dict[str, object]) -> list[dict[str, object]]:
	link = target.get("link")
	if not isinstance(link, dict):
		raise TypeError("CMake File API reply contains no linker information")
	fragments = link.get("commandFragments")
	if not isinstance(fragments, list):
		raise TypeError("CMake File API reply contains no link command fragments")
	return [item for item in fragments if isinstance(item, dict)]


def split_link_fragments(
	fragments: list[dict[str, object]], build_dir: Path
) -> tuple[list[Path], list[str]]:
	objects: list[Path] = []
	flags: list[str] = []
	for item in fragments:
		fragment = item.get("fragment")
		if not isinstance(fragment, str):
			continue
		for raw_argument in split_fragment(fragment):
			argument = raw_argument.replace("\\", "/")
			if is_object_file(argument):
				path = Path(argument)
				objects.append(path if path.is_absolute() else build_dir / path)
			elif "--out-implib" not in argument:
				path = Path(argument)
				candidate = build_dir / path
				if (
					not argument.startswith("-")
					and not path.is_absolute()
					and candidate.exists()
				):
					flags.append(normalized_path(candidate))
				else:
					flags.append(argument)
	return objects, flags


def target_source_objects(target: dict[str, object], build_dir: Path) -> list[Path]:
	sources = target.get("sources")
	if not isinstance(sources, list):
		return []
	objects: list[Path] = []
	for source in sources:
		if not isinstance(source, dict):
			continue
		path_value = source.get("path")
		if not isinstance(path_value, str) or not is_object_file(path_value):
			continue
		path = Path(path_value)
		objects.append(path if path.is_absolute() else build_dir / path)
	return objects


def target_generated_objects(build_dir: Path) -> list[Path]:
	target_object_dir = build_dir / "CMakeFiles" / f"{PROBE_TARGET}.dir"
	return [
		path
		for path in target_object_dir.rglob("*")
		if path.is_file()
		and is_object_file(path.name)
		and path.name.lower() not in {"main.cpp.obj", "main.cpp.o"}
	]


def unique_existing_objects(objects: list[Path]) -> list[Path]:
	unique_objects: list[Path] = []
	seen: set[str] = set()
	for path in objects:
		resolved = path.resolve()
		key = os.path.normcase(str(resolved))
		if key in seen:
			continue
		if not resolved.is_file():
			raise FileNotFoundError(resolved)
		seen.add(key)
		unique_objects.append(resolved)
	return unique_objects


def compile_group(target: dict[str, object]) -> dict[str, object]:
	sources = target.get("sources")
	groups = target.get("compileGroups")
	if not isinstance(sources, list) or not isinstance(groups, list):
		raise TypeError("CMake File API reply contains no compile interface")

	group_index: int | None = None
	for source in sources:
		if not isinstance(source, dict):
			continue
		path = source.get("path")
		index = source.get("compileGroupIndex")
		if isinstance(path, str) and Path(path).name == PROBE_SOURCE:
			if isinstance(index, int):
				group_index = index
			break
	if group_index is None or not 0 <= group_index < len(groups):
		raise RuntimeError(f"Could not find {PROBE_SOURCE}'s CMake compile group")

	group = groups[group_index]
	if not isinstance(group, dict) or group.get("language") != "CXX":
		raise TypeError(f"{PROBE_SOURCE}'s CMake compile group is not C++")
	return group


def dictionary_items(value: object) -> list[dict[str, object]]:
	if not isinstance(value, list):
		return []
	return [item for item in value if isinstance(item, dict)]


def compile_definitions(group: dict[str, object]) -> set[str]:
	definitions: set[str] = set()
	for item in dictionary_items(group.get("defines")):
		definition = item.get("define")
		if isinstance(definition, str):
			definitions.add(definition)
	return definitions


def compile_includes(group: dict[str, object], build_dir: Path) -> list[str]:
	flags: list[str] = []
	for item in dictionary_items(group.get("includes")):
		path_value = item.get("path")
		if not isinstance(path_value, str):
			continue
		path = Path(path_value.replace("\\", "/"))
		if not path.is_absolute():
			path = build_dir / path
		flags.extend(
			("-isystem" if item.get("isSystem") is True else "-I", normalized_path(path))
		)
	return flags


def compile_command_flags(group: dict[str, object]) -> list[str]:
	flags: list[str] = []
	for item in dictionary_items(group.get("compileCommandFragments")):
		fragment = item.get("fragment")
		if not isinstance(fragment, str):
			continue
		# Qt's own warning policy is not part of its ABI and makes generated
		# Miqt sources unnecessarily fragile across compiler releases. Likewise,
		# do not copy the probe's Release optimization into Miqt's very large
		# generated translation units; -O3 can exhaust a standard CI runner.
		flags.extend(
			flag.replace("\\", "/")
			for flag in split_fragment(fragment)
			if not flag.startswith(("-Werror", "-Wno-error", "-O", "-g"))
		)
	return flags


def compile_interface(
	target: dict[str, object], build_dir: Path
) -> tuple[list[str], list[str]]:
	group = compile_group(target)
	definitions = compile_definitions(group)
	required_definitions = {
		"QT_CORE_LIB",
		"QT_GUI_LIB",
		"QT_MULTIMEDIAWIDGETS_LIB",
		"QT_MULTIMEDIA_LIB",
		"QT_NETWORK_LIB",
		"QT_WIDGETS_LIB",
	}
	if missing := sorted(required_definitions - definitions):
		raise RuntimeError(
			"Qt's CMake compile interface is missing definitions: " + ", ".join(missing)
		)

	cpp_flags = [*(f"-D{definition}" for definition in sorted(definitions))]
	cpp_flags.extend(compile_includes(group, build_dir))
	cxx_flags = compile_command_flags(group)
	cxx_flags[:0] = ["-O1", "-g"]
	if not any(flag.startswith("-std=") for flag in cxx_flags):
		standard = group.get("languageStandard")
		version = standard.get("standard") if isinstance(standard, dict) else None
		cxx_flags.append(f"-std=c++{version if isinstance(version, str) else '17'}")
	if "-Wa,-mbig-obj" not in cxx_flags:
		cxx_flags.append("-Wa,-mbig-obj")
	return cpp_flags, cxx_flags


def link_interface(build_dir: Path) -> tuple[list[Path], list[str]]:
	target = target_reply(build_dir)
	objects, flags = split_link_fragments(command_fragments(target), build_dir)
	objects.extend(target_source_objects(target, build_dir))
	objects.extend(target_generated_objects(build_dir))
	objects = unique_existing_objects(objects)

	if not objects:
		raise RuntimeError("Qt's CMake link interface contained no generated objects")
	if import_archives := [flag for flag in flags if ".dll.a" in flag.casefold()]:
		raise RuntimeError(
			"CMake selected dynamic import libraries instead of static archives: "
			+ ", ".join(import_archives)
		)
	return objects, flags


def verify_plugin_initializers(objects: list[Path]) -> None:
	names = {path.name.casefold() for path in objects}
	missing = [
		plugin
		for plugin in REQUIRED_PLUGINS
		if f"{plugin}_init.cpp.obj".casefold() not in names
		and f"{plugin}_init.cpp.o".casefold() not in names
	]
	if missing:
		raise RuntimeError(
			"Qt's CMake interface is missing plugin initializers: " + ", ".join(missing)
		)


def write_linker_response(path: Path, arguments: list[str]) -> None:
	# GCC response files accept shell-style quoting. Keeping the complete static
	# interface here avoids Windows' command-length limit and makes Go pass it to
	# the final external link exactly once.
	path.write_text(shlex.join(arguments) + "\n", encoding="utf-8")


def environment_flags(flags: list[str]) -> str:
	return shlex.join(flag.replace("\\", "/") for flag in flags)


def main() -> None:
	environment = os.environ.copy()
	qt_version = environment.get("QT_VERSION")
	if not qt_version:
		raise RuntimeError("QT_VERSION must be set")
	gcc = required_tool("gcc")
	toolchain_bin = gcc.parent
	gxx = required_tool("g++", directory=toolchain_bin)
	cmake = required_tool("cmake", directory=toolchain_bin)
	true_tool = required_tool("true")
	mingw_prefix = toolchain_bin.parent
	qt_prefix = mingw_prefix / "qt6-static"
	if not (qt_prefix / "lib/cmake/Qt6/Qt6Config.cmake").is_file():
		raise RuntimeError(f"Static Qt is not installed at {qt_prefix}")

	environment.update(
		{
			"PKG_CONFIG_ARGN": "--static",
			"PKG_CONFIG_PATH": normalized_path(mingw_prefix / "lib/pkgconfig"),
		}
	)

	with tempfile.TemporaryDirectory(prefix="ayandict-qt-link-probe.") as temporary:
		probe_dir = Path(temporary)
		build_dir = probe_dir / "build"
		create_file_api_query(build_dir)
		run(
			[
				normalized_path(cmake),
				"-S",
				normalized_path(WINDOWS_CI_DIR / "qt-link-probe"),
				"-B",
				normalized_path(build_dir),
				"-G",
				"Ninja",
				"-DCMAKE_BUILD_TYPE=Release",
				f"-DCMAKE_CXX_COMPILER={normalized_path(gxx)}",
				"-DCMAKE_FIND_LIBRARY_SUFFIXES=.a",
				"-DCMAKE_DISABLE_FIND_PACKAGE_harfbuzz=ON",
				"-DCMAKE_DISABLE_FIND_PACKAGE_PCRE2=ON",
				"-DCMAKE_EXE_LINKER_FLAGS=-static -static-libgcc -static-libstdc++",
				"-DOPENSSL_USE_STATIC_LIBS=ON",
				"-DPKG_CONFIG_ARGN=--static",
				"-DZLIB_USE_STATIC_LIBS=ON",
				f"-DCMAKE_PREFIX_PATH={normalized_path(qt_prefix)}",
				f"-DAYANDICT_QT_VERSION={qt_version}",
			],
			env=environment,
		)
		run(
			[
				normalized_path(cmake),
				"--build",
				normalized_path(build_dir),
				"--verbose",
			],
			env=environment,
		)
		probe_executable = build_dir / f"{PROBE_TARGET}.exe"
		run(
			[
				normalized_path(Path(sys.executable)),
				normalized_path(WINDOWS_CI_DIR / "verify-binary.py"),
				"--imports-only",
				normalized_path(probe_executable),
			],
			env=environment,
		)
		run(
			[normalized_path(probe_executable)],
			env=environment,
		)

		target = target_reply(build_dir)
		cpp_flags, cxx_flags = compile_interface(target, build_dir)
		objects, link_flags = link_interface(build_dir)
		verify_plugin_initializers(objects)
		linker_response = probe_dir / "static-link.rsp"
		write_linker_response(
			linker_response,
			[
				"-static",
				"-static-libgcc",
				"-static-libstdc++",
				*(normalized_path(path) for path in objects),
				*link_flags,
			],
		)

		environment.update(
			{
				"CC": normalized_path(gcc),
				"CGO_CFLAGS": "-O1 -g -Wa,-mbig-obj",
				"CGO_CPPFLAGS": environment_flags(cpp_flags),
				"CGO_CXXFLAGS": environment_flags(cxx_flags),
				"CGO_ENABLED": "1",
				# The package-level CGo probe links may fail without Qt; cmd/go
				# records that fact and uses the explicitly requested external link.
				"CGO_LDFLAGS": "",
				"CXX": normalized_path(gxx),
				"GOARCH": "amd64",
				"GOOS": "windows",
				# Miqt generates unusually large C/C++ translation units. Limit both
				# Go package compilation and compiler-process concurrency to fit the
				# memory available on GitHub's standard Windows runners.
				"GOMAXPROCS": "2",
				"PKG_CONFIG": normalized_path(true_tool),
			}
		)
		run(
			[
				"go",
				"build",
				"-p=1",
				"-v",
				"-trimpath",
				"-ldflags",
				(
					"-w -H=windowsgui -linkmode=external "
					f"-extldflags=@{normalized_path(linker_response)}"
				),
				"-o",
				"ayandict-unstripped.exe",
			],
			env=environment,
			cwd=REPO_ROOT,
		)


if __name__ == "__main__":
	main()
