from __future__ import annotations

# ruff: noqa: INP001
import importlib.util
import unittest
from pathlib import Path
from typing import TYPE_CHECKING

if TYPE_CHECKING:
	import types


def load_build_go() -> types.ModuleType:
	path = Path(__file__).with_name("build-go.py")
	spec = importlib.util.spec_from_file_location("build_go", path)
	assert spec is not None
	assert spec.loader is not None
	module = importlib.util.module_from_spec(spec)
	spec.loader.exec_module(module)
	return module


class SplitFragmentTest(unittest.TestCase):
	def test_preserves_windows_paths(self) -> None:
		build_go = load_build_go()
		path = (
			r"D:\a\_temp\msys64\ucrt64\qt6-static\lib\objects-Release"
			r"\Widgets_resources_1\.qt\rcc\qrc_qstyle_init.cpp.obj"
		)
		self.assertEqual(build_go.split_fragment(path), [path])

	def test_removes_outer_quotes_only(self) -> None:
		build_go = load_build_go()
		path = r"D:\a path\Qt\plugin_init.cpp.obj"
		self.assertEqual(
			build_go.split_fragment(f'"{path}" -lQt6Core'),
			[path, "-lQt6Core"],
		)


class CompileFlagsTest(unittest.TestCase):
	def test_drops_probe_optimization_and_debug_flags(self) -> None:
		build_go = load_build_go()
		group = {
			"compileCommandFragments": [
				{"fragment": "-O3 -g -DNDEBUG -ffunction-sections"}
			]
		}
		self.assertEqual(
			build_go.compile_command_flags(group),
			["-DNDEBUG", "-ffunction-sections"],
		)


if __name__ == "__main__":
	unittest.main()
