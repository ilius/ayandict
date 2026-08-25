# Static Windows build

The Windows CI job produces one standalone `ayandict.exe` for x86-64 Windows 10
1809 or newer and Windows 11. Qt, its selected plugins, the MinGW C/C++
runtimes, and required third-party libraries are linked statically. Only
Windows system DLL imports are allowed by `verify-binary.py`.

The build pins the MSYS2 UCRT64 static Qt package at Qt 6.11.2-1 and verifies
its SHA-256 before installation. Qt's CMake metadata is the source of truth for
the complete ordered static link interface. `build-go.py` builds a small CMake
probe, collects Qt's generated resource/plugin initializer objects, and passes
that interface directly to the final external Go/CGo link.

Miqt's `windowsqtstatic` build tag is intentionally omitted. CMake's generated
plugin initializer objects provide the static imports, including multimedia;
using both mechanisms would register the platform and style plugins twice.

The executable includes these Qt plugins:

- Windows platform integration and native style
- Windows Media Foundation multimedia
- GIF, ICO, JPEG, SVG, WebP, and TIFF image formats

This build uses Qt's native Windows Media Foundation backend and ships no
FFmpeg. Format and codec availability is determined by the target Windows
installation; applications can query `QMediaFormat` at runtime. Qt has marked
this backend deprecated since Qt 6.10.

## Pinned Qt inputs and release obligations

- [MSYS2 static Qt package](https://packages.msys2.org/packages/mingw-w64-ucrt-x86_64-qt6-static)
- [Exact MSYS2 package source](https://mirror.msys2.org/mingw/sources/mingw-w64-qt6-static-6.11.2-1.src.tar.zst)
- [Exact Qt 6.11.2 source](https://download.qt.io/official_releases/qt/6.11/6.11.2/single/qt-everywhere-src-6.11.2.tar.xz)

The Qt binary package is pinned and integrity-checked. The GitHub-hosted runner,
actions, compiler, build tools, and third-party packages are rolling inputs, so
the complete build environment is not bit-for-bit reproducible.

AyanDict is AGPL-3.0-or-later. Distributing a statically linked Qt build also
requires satisfying the applicable Qt and third-party notice, corresponding
source, and relinking obligations. The executable can remain standalone, but a
release should accompany it with the required notices, license texts, and an
exact corresponding-source/relink offer or bundle.
