# Development instructions

## Build and test constraints

- Do not automatically build or test the whole application. It pulls in MiQt, whose first build is long-running and RAM-intensive.
- Build the complete application only when the user explicitly requests it.
- When a Qt-related change needs a build, use `./build-agent`. This avoids an unnecessary MiQt rebuild by copying only existing MiQt artifacts into the sandbox-writable Go build cache before building.
- Do not run `go test` on packages that have no `_test.go` files.
- When asked to commit, never run `go fmt`, `gofmt` or `go test`.
