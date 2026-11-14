# State: pigeon

Overwritten to stay current (loopers and agents). Last updated: 2026-11-14 by looper (#11).

- **Version / line:** v0.1.3 (the 0.1 line; its image is #2); tacplus floor 2.6 (#8; coop declares 2.6), the build
  pinned by `tacplus.lock` to tacplus 2.6.1 and its `build:2.6.1` image by digest
- **Works:** the generic pigeon host `pigeonhost` with the built-in kind `echo.ping`; every message a task causes (its
  Ready, result or failure) carries the task's `trace_id` (#7). In the pinned image (tacplus 2.6.1,
  Go 1.27.3), `just gates` (`tacplus build` leaving every tracked file as it was, `go mod tidy` with no diff, gofmt,
  `go build`, `go vet`, `go test ./...`; the justfile is #9's) passes; govulncheck v1.9.0 finds no vulnerabilities (#8).
- **Released / deployed:** tagged (v0.1.0, v0.1.1, v0.1.2, v0.1.3); no image published yet (#2).
- **Known gaps:** none in this repository. The looper's `tacplus` wrapper runs the pinned release, in its line's
  `build:2.6` image; the Go gates run in the lock's `build_image`. Building main in tacplus-dev/build from the pin is #9.
- **Blocked on:** #2 (the release image) on tacplus-dev/looper#115: a report's token cannot push `.github/workflows/`,
  so its `release.yml` waits; #10 (the tacplus 2.6.8 pin) on tacplus-dev/looper#113: its build job never reports a
  verdict.
