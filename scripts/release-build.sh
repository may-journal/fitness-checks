#!/bin/sh
# Builds the check a release tag names for four platforms, stripped, into a tar and gzip with checksums, plus <tag>.describe.json.
# Usage: scripts/release-build.sh <tag> <out-dir>, with archives named <tag>-<os>-<arch>.
# cspell:ignore trimpath buildvcs ldflags
set -eu
cd "$(dirname "$0")/.."

tag="$1"
out="$2"
component="${tag%-v*}"
[ -d "cmd/$component" ] || { echo "release-build: no folder cmd/$component for tag $tag" >&2; exit 1; }
rm -rf "$out"
mkdir -p "$out"

host="$(go env GOOS)/$(go env GOARCH)"
for platform in linux/amd64 linux/arm64 darwin/amd64 darwin/arm64; do
  os="${platform%/*}"
  arch="${platform#*/}"
  work="$(mktemp -d)"
  CGO_ENABLED=0 GOOS="$os" GOARCH="$arch" go build -trimpath -buildvcs=false -ldflags="-s -w" \
    -o "$work/$component" "./cmd/$component"
  # The describe line comes from the shipped program the host can run.
  if [ "$platform" = "$host" ]; then
    "$work/$component" --describe > "$out/$tag.describe.json"
  fi
  archive="$tag-$os-$arch.tar"
  COPYFILE_DISABLE=1 tar -C "$work" -cf "$out/$archive" "$component"
  gzip -9 -n -k "$out/$archive"
  (cd "$out" && shasum -a 256 "$archive" > "$archive.sha256" && shasum -a 256 "$archive.gz" > "$archive.gz.sha256")
  rm -rf "$work"
done
[ -s "$out/$tag.describe.json" ] || go run "./cmd/$component" --describe > "$out/$tag.describe.json"
# A pattern that does not compile would make the engine reject the check, so it stops the release here.
go run ./scripts/describe-patterns < "$out/$tag.describe.json"
