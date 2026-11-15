#!/usr/bin/env bash
# Assembles the release image (.github/workflows/release.yml) from the source tree's Dockerfile, once per
# platform, and publishes them as one multi-architecture image.
#
# Usage:  ci/release-image.sh SRC OUT VERSION REVISION [TAG...]
#
# The Dockerfile takes its binaries from beside it, as the fleet's image pipeline puts them there: for each
# platform in PLATFORMS (default "arm64") this stages a build context with SRC/Dockerfile,
# OUT/pigeonhost-linux-<arch> as bin/pigeonhost and OUT/launcher-linux-<arch>, and builds it labelled with VERSION
# and REVISION. With TAGs, each platform is pushed by digest and $IMAGE:<TAG> (each TAG) is pointed at all of
# them. Without, nothing is pushed: the dry run.
set -euo pipefail

src="$1"
out="$2"
version="$3"
revision="$4"
shift 4
image="${IMAGE:-ghcr.io/tacplus-dev/pigeon}"

digests=()
for arch in ${PLATFORMS:-arm64}; do
  case "$arch" in
    amd64|arm64) ;;
    *) echo "release-image.sh: unsupported platform $arch (amd64, arm64)" >&2; exit 1 ;;
  esac
  ctx="$out/image/$arch"
  rm -rf "$ctx"
  mkdir -p "$ctx/bin"
  cp "$src/Dockerfile" "$ctx/"
  cp "$out/pigeonhost-linux-$arch" "$ctx/bin/pigeonhost"
  cp "$out/launcher-linux-$arch" "$ctx/"

  output="type=cacheonly"
  (($#)) && output="type=image,name=$image,push-by-digest=true,name-canonical=true,push=true"
  docker buildx build --platform "linux/$arch" \
    --build-arg "VERSION=$version" --build-arg "VCS_REF=$revision" \
    --output "$output" --metadata-file "$ctx.json" "$ctx"
  if (($#)); then
    digest="$(jq -r '."containerimage.digest"' "$ctx.json")"
    echo "linux/$arch: $image@$digest"
    digests+=("$image@$digest")
  fi
done

(($#)) || exit 0
refs=()
for tag in "$@"; do refs+=(-t "$image:$tag"); done
docker buildx imagetools create "${refs[@]}" "${digests[@]}"
