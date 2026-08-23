# macOS CI scripts

`config.toml` is the source of truth for Qt and Go versions, the Qt source URL,
installation root, and cache revisions. The GitHub Actions workflow reads its
cache paths and tool versions from the same file used by these scripts.

The Python scripts require Python 3.11 or newer and are intended to run from the
repository checkout on a macOS runner. `build-go.py` asks Qt's CMake integration
for the complete static link interface before invoking Go, including Qt's
generated resources and static plugins.
