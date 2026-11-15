# State: pigeon

Overwritten to stay current (loopers and agents). Last updated: 2026-11-15 by looper (#2).

- **Version / line:** v0.1.3 (the 0.1 line), its image `ghcr.io/tacplus-dev/pigeon:0.1.3`, which the line tag `:0.1`
  names (linux/arm64; #2); tacplus floor 2.6 (#8; coop declares 2.6), the build pinned by `tacplus.lock` to tacplus
  2.6.1 and its `build:2.6.1` image by digest
- **Works:** the generic pigeon host `pigeonhost` with the built-in kind `echo.ping`; every message a task causes (its
  Ready, result or failure) carries the task's `trace_id` (#7). In the pinned image (tacplus 2.6.1,
  Go 1.27.3), `just gates` (`tacplus build` leaving every tracked file as it was, `go mod tidy` with no diff, gofmt,
  `go build`, `go vet`, `go test ./...`; the justfile is #9's) passes; govulncheck v1.9.0 finds no vulnerabilities (#8).
  A pushed `vX.Y.Z` tag runs `.github/workflows/release.yml` (#2; README, Releases): the gates on the tag in its
  pinned build image, then the image `X.Y.Z` (the line tag `X.Y` only for the line's highest release) and the GitHub
  release with the binaries; its scripts are `ci/` (`go test ./ci/`). The package is public: pulls need no credential.
- **Released / deployed:** v0.1.1, v0.1.2 and v0.1.3 are published (image and GitHub release, linux/arm64; v0.1.3 is
  Latest and `:0.1`). v0.1.0 is tagged with no image: its tree predates `tacplus.lock`, and the floor's `build:2.5`
  image has no `tacplus` on the PATH of the non-root user the workflow builds as, so v0.1.1 is the first image, as
  #2 allows. Not deployed.
- **Known gaps:** the image is linux/arm64 only, as the fleets run (`PLATFORMS` in release.yml turns amd64 on;
  PROPOSALS.md). The looper's `tacplus` wrapper runs the pinned release, in its line's `build:2.6` image; the Go gates
  run in the lock's `build_image`. Building main in tacplus-dev/build from the pin is #9.
- **Blocked on:** #10 (the tacplus 2.6.8 pin) on tacplus-dev/looper#113: its build job never reports a verdict.
