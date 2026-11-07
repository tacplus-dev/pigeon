# State: pigeon

Overwritten to stay current (loopers and agents). Last updated: 2026-11-07 by looper (#8).

- **Version / line:** v0.1.1 (the 0.1 line; its image is #2); tacplus floor 2.6 (#8; coop declares 2.6), the build
  pinned by `tacplus.lock` to tacplus 2.6.1 and its `build:2.6.1` image by digest
- **Works:** the generic pigeon host `pigeonhost` with the built-in kind `echo.ping`; every message a task causes (its
  Ready, result or failure) carries the task's `trace_id` (#7). In the pinned image (tacplus 2.6.1,
  Go 1.27.3), `tacplus build`, `go mod tidy` (no diff), gofmt, `go vet`, `go build` and `go test ./...` pass;
  govulncheck v1.9.0 finds no vulnerabilities (#8).
- **Released / deployed:** tagged (v0.1.0, v0.1.1); no image published yet (#2).
- **Known gaps:** none in this repository. The looper's `tacplus` wrapper runs the pinned release, in its line's
  `build:2.6` image; the Go gates run in the lock's `build_image`. Building main in tacplus-dev/build from the pin is #9.
- **Blocked on:** #2 on tacplus-dev/looper#12.
