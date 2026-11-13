#!/usr/bin/env bash
# Runs a command in the build image a source tree builds in (ci/build-image.sh), on that tree, as the caller's
# user, with no credential. The tree must be in this checkout (the release checks the tag out into src/), which
# is mounted at its own path, so the command may be one of this checkout's ci/ scripts.
#
# Usage:  ci/in-build-image.sh DIR COMMAND [ARG...]
set -euo pipefail

root="$(cd "$(dirname "$0")/.." && pwd)"
dir="$(cd "$1" && pwd)"
shift
[[ "$dir" == "$root" || "$dir" == "$root"/* ]] || { echo "in-build-image.sh: $dir is not in $root" >&2; exit 1; }
image="$("$root/ci/build-image.sh" "$dir")"
echo "in-build-image.sh: $image" >&2

docker run --rm -u "$(id -u):$(id -g)" -e HOME=/tmp -v "$root:$root" -w "$dir" "$image" "$@"
