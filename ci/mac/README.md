# macOS CI scripts

`config.toml` is the source of truth for Qt and Go versions, the Qt source URL,
installation root, and cache revisions. The GitHub Actions workflow reads its
cache paths and tool versions from the same file used by these scripts.

The scripts are intended to run from the repository checkout on a macOS runner.
`build-go.sh` asks Qt's CMake integration for the complete static link interface
before invoking Go, including Qt's generated resources and static plugins.
