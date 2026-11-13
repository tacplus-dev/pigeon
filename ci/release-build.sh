#!/usr/bin/env bash
# The release build of a source tree, run in its build image (ci/in-build-image.sh) at the tree's root: its
# gates, then the static pigeonhost binary for each platform, as OUT/pigeonhost-linux-<arch>.
#
# Usage:  ci/release-build.sh OUT PLATFORM...
#
# The gates are the tree's own `just gates` (the justfile). A tree from before the justfile (v0.1.0, v0.1.1)
# runs the gates that justfile wrote down, the ones its main passed: tacplus build leaving every tracked file as
# it was, go mod tidy with no diff, gofmt, go build, go vet, go test.
set -euo pipefail

out="$1"
shift
mkdir -p "$out"

if [[ -f justfile ]] && just --summary | tr ' ' '\n' | grep -qx gates; then
  just gates
else
  echo "release-build.sh: no justfile gates here: running the gates it wrote down" >&2
  before="$(git diff HEAD | sha256sum)"
  tacplus build
  [[ "$(git diff HEAD | sha256sum)" == "$before" ]] || { echo "tacplus build changed tracked files" >&2; git status --short >&2; exit 1; }
  go mod tidy
  git diff --exit-code -- go.mod go.sum || { echo "go.mod/go.sum are not tidy" >&2; exit 1; }
  unformatted="$(gofmt -l cmd internal)"
  [[ -z "$unformatted" ]] || { echo "gofmt -l reports: $unformatted" >&2; exit 1; }
  go build ./...
  go vet ./...
  go test ./...
fi

for arch in "$@"; do
  CGO_ENABLED=0 GOOS=linux GOARCH="$arch" go build -trimpath -ldflags="-s -w" -o "$out/pigeonhost-linux-$arch" ./cmd/pigeonhost
done
go version
