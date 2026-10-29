# State: pigeon

Overwritten to stay current (loopers and agents). Last updated: 2026-10-29 by looper (#5).

- **Version / line:** unreleased (no tag; v0.1.0 is #2); tacplus floor 2.5 (#5; stratum, via coop, needs 2.5)
- **Works:** the generic pigeon host `pigeonhost` with the built-in kind `echo.ping`. On the 2.5 line (tacplus 2.5.1
  plus Go 1.26.10), `tacplus build`, `go mod tidy` (no diff), gofmt, `go vet`, `go build` and `go test ./...` pass.
- **Released / deployed:** nowhere yet.
- **Known gaps:** the gates fail through the looper's `tacplus` wrapper: it runs the deprecated compiler image
  `tacplus/tacplus:2.5`, which has no `go`, so Go generation fails at goimports (tacplus-dev/looper#27).
- **Blocked on:** tacplus-dev/looper#27 (gates in the looper); #2 also on tacplus-dev/looper#12.
