#!/usr/bin/env bash
# How the GitHub release of a version is marked (.github/workflows/release.yml): the flag for `gh release create`.
#
# Usage:  ci/release-latest.sh vX.Y.Z[-pre]
#
# Run in the checkout with its tags fetched. Prints --latest when vX.Y.Z is the highest release of all, by
# version; --latest=false otherwise, so a patch on an older line (or an older release published later) never
# becomes Latest; --prerelease for a pre-release. Tested by ci/release_test.go.
set -euo pipefail

ref="$1"
if [[ "$ref" =~ ^v[0-9]+\.[0-9]+\.[0-9]+-[0-9A-Za-z.-]+$ ]]; then
  echo --prerelease
  exit 0
fi
if [[ ! "$ref" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo "release-latest.sh: $ref is not vMAJOR.MINOR.PATCH" >&2
  exit 1
fi
highest="$(git tag -l 'v*' | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+$' | sort -t. -k1.2,1n -k2,2n -k3,3n | tail -1 || true)"
if [[ "$highest" == "$ref" ]]; then
  echo --latest
else
  echo --latest=false
fi
