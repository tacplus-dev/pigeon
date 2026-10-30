# State: pigeon

Overwritten to stay current (loopers and agents). Last updated: 2026-10-30 by looper (#5).

- **Version / line:** unreleased (no tag; v0.1.0 is #2); tacplus floor 2.5 (#5; stratum, via coop, needs 2.5)
- **Works:** the generic pigeon host `pigeonhost` with the built-in kind `echo.ping`. On the 2.5 line (tacplus 2.5.3
  plus Go 1.26.10), `tacplus build`, `go mod tidy` (no diff), gofmt, `go vet`, `go build` and `go test ./...` pass;
  govulncheck v1.9.0 finds no vulnerabilities (golang.org/x/net v0.62.0, #6).
- **Released / deployed:** nowhere yet.
- **Known gaps:** none in this repository. The gates pass through the looper's fixed `tacplus` wrapper
  (tacplus-dev/looper@71591c7, #27: tacplus 2.5.3 in `build:2.5`, 2026-10-30); a looper image built before that
  fix still runs the Go-less compiler image, so Go generation fails there at goimports.
- **Blocked on:** #2 on tacplus-dev/looper#12.
