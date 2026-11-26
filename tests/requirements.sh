#!/bin/sh
# Runs the requirement suite the way a user meets these checks. The installer
# unpacks the pinned engine with this repo's check programs beside it.
set -eu
cd "$(dirname "$0")/.."

# The installer takes only X.Y.Z. A pseudo-version such as v1.61.1-0.<time>-<sha>
# builds on release 1.61.0, so the engine reports that release.
version="$(go list -m -f "{{.Version}}" github.com/may-journal/fitness-runner/go)"
version="${version#v}"
case "$version" in
*-0.*)
  base="${version%%-*}"
  version="${base%.*}.$((${base##*.} - 1))"
  ;;
esac
platform="$(go env GOOS)-$(go env GOARCH)"
top="fitness-$version-$platform"
work="$PWD/out/requirements"
rm -rf "$work"
mkdir -p "$work/$top"

go build -ldflags "-s -w -X main.version=$version" -o "$work/$top/" \
  github.com/may-journal/fitness-runner/go/cmd/fitness ./cmd/...
COPYFILE_DISABLE=1 tar -C "$work" -czf "$work/bundle.tar.gz" "$top"
digest="$(shasum -a 256 "$work/bundle.tar.gz" | cut -d' ' -f1)"

# The installer finds the bundle already in its cache, under the hash its
# build pins, so nothing downloads.
FITNESS_CACHE_DIR="$work/cache"
mkdir -p "$FITNESS_CACHE_DIR/$top-$digest"
mv "$work/bundle.tar.gz" "$FITNESS_CACHE_DIR/$top-$digest/archive.tar.gz"
go build -ldflags "-X main.version=$version -X main.bundleHashes=$platform=$digest" \
  -o "$work/fitness-install" github.com/may-journal/fitness-runner/go/cmd/fitness-install
FITNESS_INSTALL="$work/fitness-install"
export FITNESS_CACHE_DIR FITNESS_INSTALL

# The repo root's requirements folder and each package's own, such as
# cmd/<check>/requirements, hold the requirement tests.
folders="$(go list ./... | grep '/requirements$')"
[ -n "$folders" ] || { echo "no requirements folder found" >&2; exit 1; }
# internal/acceptance's TakeTurn caps how many of these tests run at once
# across the whole machine.
# shellcheck disable=SC2086
exec go test $folders -count=1 "$@"
