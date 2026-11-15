#!/usr/bin/env bash
# The build image a source tree builds in (ci/in-build-image.sh): the image its tacplus.lock pins by digest
# (`tacplus pin` writes it). A tree from before the lock (v0.1.0) builds in the build image of the floor its
# tacplus.yaml declares, `build:X.Y`, which is that line's latest patch: not reproducible, and said so.
#
# Usage:  ci/build-image.sh DIR
#
# Tested by ci/release_test.go.
set -euo pipefail

dir="$1"
if [[ -f "$dir/tacplus.lock" ]]; then
  image="$(sed -n -E 's/^build_image:[[:space:]]*"?([^"[:space:]]+)"?.*$/\1/p' "$dir/tacplus.lock" | head -1)"
  [[ "$image" == *@sha256:* ]] || { echo "build-image.sh: $dir/tacplus.lock names no build_image by digest" >&2; exit 1; }
  echo "$image"
  exit 0
fi
floor="$(sed -n -E 's/^tacplus:[[:space:]]*"?([0-9]+\.[0-9]+)"?.*$/\1/p' "$dir/tacplus.yaml" 2>/dev/null | head -1)"
[[ -n "$floor" ]] || { echo "build-image.sh: $dir has no tacplus.lock and no tacplus floor in tacplus.yaml" >&2; exit 1; }
echo "build-image.sh: $dir has no tacplus.lock: building in the floor's latest build image, build:$floor (not pinned)" >&2
echo "ghcr.io/tacplus-dev/tacplus/build:$floor"
