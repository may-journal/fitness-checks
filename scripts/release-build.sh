#!/bin/sh
# Builds the check a release tag names for four platforms into one archive and checksum file each, named <tag>-<os>-<arch>.
# Usage: scripts/release-build.sh <tag> <out-dir>
# cspell:ignore trimpath buildvcs
set -eu
cd "$(dirname "$0")/.."

tag="$1"
out="$2"
component="${tag%-v*}"
[ -d "cmd/$component" ] || { echo "release-build: no folder cmd/$component for tag $tag" >&2; exit 1; }
rm -rf "$out"
mkdir -p "$out"

for platform in linux/amd64 linux/arm64 darwin/amd64 darwin/arm64; do
  os="${platform%/*}"
  arch="${platform#*/}"
  work="$(mktemp -d)"
  CGO_ENABLED=0 GOOS="$os" GOARCH="$arch" go build -trimpath -buildvcs=false \
    -o "$work/$component" "./cmd/$component"
  archive="$tag-$os-$arch.tar"
  COPYFILE_DISABLE=1 tar -C "$work" -cf "$out/$archive" "$component"
  (cd "$out" && shasum -a 256 "$archive" > "$archive.sha256")
  rm -rf "$work"
done
