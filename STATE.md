# State: pigeon

Overwritten to stay current (loopers and agents). Last updated: 2026-11-24 by looper (#17).

- **Version / line:** v0.1.8 (the 0.1 line), its image `ghcr.io/tacplus-dev/pigeon:0.1.8`, which the line tag `:0.1`
  names (linux/arm64; #2); tacplus floor 2.6 (#8; coop declares 2.6), the build pinned by `tacplus.lock` to tacplus
  2.6.14 and its `build:2.6.14` image by digest (#16; 2.6.13 was #10's, Go 1.27.4, fixing GO-2026-6859)
- **Works:** the generic pigeon host `pigeonhost` with the built-in kind `echo.ping`; every message a task causes (its
  Ready, result or failure) carries the task's `trace_id` (#7). In the pinned image (tacplus 2.6.14), `just gates` (`tacplus build` leaving every tracked file as it was, `go mod tidy` with no diff, gofmt,
  `go build`, `go vet`, `go test ./...`; the justfile is #9's) passes; govulncheck v1.9.0 finds no vulnerabilities in the
  `pigeonhost` it builds (#8, #10).
  A pushed `vX.Y.Z` tag runs `.github/workflows/release.yml` (#2; README, Releases): the gates on the tag in its
  pinned build image, then the image `X.Y.Z` (the line tag `X.Y` only for the line's highest release) and the GitHub
  release with the binaries; its scripts are `ci/` (`go test ./ci/`). The package is public: pulls need no credential.
  The image runs `apk upgrade --no-cache` before its `apk add`, so each release carries its Alpine line's
  security fixes (zlib 1.4-r0 for CVE-2026-4738, which v0.1.4's image lacks; #13, shipped in v0.1.5).
  The image's launcher is built from launcher v0.1.0 (`LAUNCHER_COMMIT` 0e139ca, go1.27.4, fixing GO-2026-6859, which
  the launcher of v0.1.5 and earlier images has; its health probe verifies the child's certificate; #12, shipped in
  v0.1.6). v0.1.7 is v0.1.6's tree rebuilt for openssl 3.5.9-r1 (libssl3, libcrypto3: CVE-2026-107496, -88822,
  -88823, -95691, -95692, which v0.1.6's image has; #15): trivy 0.76.1 finds no vulnerability in its image.
  v0.1.8 is v0.1.7's tree built and gated on the tacplus 2.6.14 pin (#16).
  On main (unreleased; v0.1.9 is its patch, #17) the image's base is pinned by digest, `alpine:3.24.4@sha256:f59abe65…`
  (`BASE_IMAGE`; `TestDockerfileBasesPinnedByDigest` fails on any `FROM` without one), so a rebuild ships the base it
  was gated on; the nightly maintenance moves the digest (tacplus-dev/looper#234).
- **Released / deployed:** v0.1.1 to v0.1.3 are published (image and GitHub release, linux/arm64); v0.1.4 (#10, the
  tacplus 2.6.13 pin) and v0.1.5 (#13) are published; v0.1.6 (#12, launcher v0.1.0) is published; v0.1.7
  (#15, openssl 3.5.9-r1) is published; v0.1.8 (#16, the tacplus 2.6.14 pin) is published as Latest and `:0.1`. v0.1.0 is tagged with no image: its tree predates
  `tacplus.lock`, and the floor's `build:2.5` image has no `tacplus` on the PATH of the non-root user the workflow
  builds as, so v0.1.1 is the first image, as #2 allows. Not deployed.
- **Known gaps:** the image is linux/arm64 only, as the fleets run (`PLATFORMS` in release.yml turns amd64 on;
  PROPOSALS.md). The looper's `tacplus` wrapper runs the pinned release, in its line's `build:2.6` image; the Go gates
  run in the lock's `build_image`. Building main in tacplus-dev/build from the pin is #9.
- **Blocked on:** nothing.
