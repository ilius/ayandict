#!/usr/bin/env python3

from __future__ import annotations

import os
import re
import shutil
import subprocess
import sys
from pathlib import Path

DLL_PATTERN = re.compile(r"^\s*DLL Name:\s*(\S+)", re.MULTILINE)
WINDOWS_SYSTEM_DLLS = frozenset(
	{
		"advapi32.dll",
		"authz.dll",
		"avrt.dll",
		"bcrypt.dll",
		"cfgmgr32.dll",
		"combase.dll",
		"comctl32.dll",
		"comdlg32.dll",
		"crypt32.dll",
		"d2d1.dll",
		"d3d9.dll",
		"d3d11.dll",
		"d3d12.dll",
		"dbghelp.dll",
		"dcomp.dll",
		"dnsapi.dll",
		"dwmapi.dll",
		"dwrite.dll",
		"dxgi.dll",
		"dxva2.dll",
		"evr.dll",
		"gdi32.dll",
		"hid.dll",
		"imagehlp.dll",
		"imm32.dll",
		"iphlpapi.dll",
		"kernel32.dll",
		"kernelbase.dll",
		"mf.dll",
		"mfcore.dll",
		"mfplat.dll",
		"mfreadwrite.dll",
		"mpr.dll",
		"msvcrt.dll",
		"ncrypt.dll",
		"netapi32.dll",
		"normaliz.dll",
		"ntdll.dll",
		"ole32.dll",
		"oleacc.dll",
		"oleaut32.dll",
		"opengl32.dll",
		"powrprof.dll",
		"propsys.dll",
		"psapi.dll",
		"rpcrt4.dll",
		"runtimeobject.dll",
		"secur32.dll",
		"setupapi.dll",
		"shcore.dll",
		"shell32.dll",
		"shlwapi.dll",
		"synchronization.dll",
		"ucrtbase.dll",
		"urlmon.dll",
		"user32.dll",
		"userenv.dll",
		"uxtheme.dll",
		"usp10.dll",
		"version.dll",
		"winhttp.dll",
		"wininet.dll",
		"winmm.dll",
		"winspool.drv",
		"wintrust.dll",
		"wldap32.dll",
		"ws2_32.dll",
		"wtsapi32.dll",
	}
)
SYSTEM_DLL_PREFIXES = ("api-ms-win-", "ext-ms-win-")
FORBIDDEN_DLL_PREFIXES = (
	"libb2",
	"libbrotli",
	"libbz2",
	"libcrypto",
	"libdeflate",
	"libfreetype",
	"libgcc_s",
	"libglib",
	"libgraphite",
	"libharfbuzz",
	"libiconv",
	"libintl",
	"libjbig",
	"libjpeg",
	"liblerc",
	"liblzma",
	"libpcre",
	"libpng",
	"libsharpyuv",
	"libssl",
	"libstdc++",
	"libtiff",
	"libvulkan",
	"libwebp",
	"libwinpthread",
	"libzstd",
	"msys-",
	"qt6",
	"vulkan-1",
	"zlib",
)
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


def required_tool(name: str) -> str:
	path = shutil.which(name)
	if not path:
		raise RuntimeError(f"Required tool is not installed: {name}")
	return path


def imported_dlls(binary: Path) -> list[str]:
	output = subprocess.run(
		[required_tool("objdump"), "-p", str(binary)],
		check=True,
		stdout=subprocess.PIPE,
		text=True,
	).stdout
	return sorted(set(DLL_PATTERN.findall(output)), key=str.casefold)


def is_windows_system_dll(name: str) -> bool:
	folded = name.casefold()
	if folded.startswith(FORBIDDEN_DLL_PREFIXES):
		return False
	if folded in WINDOWS_SYSTEM_DLLS or folded.startswith(SYSTEM_DLL_PREFIXES):
		return True
	windows_directory = os.environ.get("WINDIR")
	return bool(
		windows_directory and (Path(windows_directory) / "System32" / name).is_file()
	)


def verify_imports(binary: Path) -> None:
	imports = imported_dlls(binary)
	if not imports:
		raise RuntimeError("PE import table contained no Windows system DLLs")
	unexpected = [name for name in imports if not is_windows_system_dll(name)]
	if unexpected:
		details = "\n".join(f"  {name}" for name in unexpected)
		raise RuntimeError(f"Executable has non-system DLL dependencies:\n{details}")
	sys.stdout.write("Windows system DLL imports:\n")
	sys.stdout.write("\n".join(f"  {name}" for name in imports) + "\n")


def verify_plugins(binary: Path) -> None:
	symbols = subprocess.run(
		[required_tool("nm"), "-C", "--defined-only", str(binary)],
		check=True,
		stdout=subprocess.PIPE,
		text=True,
	).stdout
	missing = [
		name for name in REQUIRED_PLUGINS if f"qt_static_plugin_{name}()" not in symbols
	]
	if missing:
		details = "\n".join(f"  {name}" for name in missing)
		raise RuntimeError(f"Executable is missing static Qt plugins:\n{details}")
	sys.stdout.write("Verified static Qt plugins:\n")
	sys.stdout.write("\n".join(f"  {name}" for name in REQUIRED_PLUGINS) + "\n")


def main() -> None:
	arguments = sys.argv[1:]
	imports_only = False
	if arguments and arguments[0] == "--imports-only":
		imports_only = True
		arguments = arguments[1:]
	if len(arguments) != 1:
		raise SystemExit(f"usage: {Path(sys.argv[0]).name} [--imports-only] BINARY")

	binary = Path(arguments[0]).resolve()
	if not binary.is_file():
		raise FileNotFoundError(binary)
	verify_imports(binary)
	if not imports_only:
		verify_plugins(binary)


if __name__ == "__main__":
	main()
