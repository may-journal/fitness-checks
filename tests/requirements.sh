#!/bin/sh
# Runs the requirement suite the way a user meets these checks. The installer
# unpacks the pinned engine with this repo's check programs beside it.
# REQUIREMENT_SHARD=<i>/<n> runs share i of n of the root requirements package,
# and REQUIREMENT_SHARD=checks runs the folders of the checks REQUIREMENT_CHECKS
# names as a JSON array, so CI splits the suite over parallel jobs.
set -eu
cd "$(dirname "$0")/.."

# The repo root's requirements folder and each package's own, such as
# cmd/<check>/requirements, hold the requirement tests.
case "${REQUIREMENT_SHARD:-}" in
"") folders="$(go list ./... | grep '/requirements$')" ;;
checks) folders="$(scripts/requirement-shard.sh checks "${REQUIREMENT_CHECKS:-}")" ;;
*)
  folders=./requirements
  pattern="$(scripts/requirement-shard.sh "$REQUIREMENT_SHARD")"
  set -- -run "$pattern" "$@"
  ;;
esac
if [ -z "$folders" ]; then
  [ -n "${REQUIREMENT_SHARD:-}" ] || { echo "no requirements folder found" >&2; exit 1; }
  echo "requirements.sh: the changed checks own no requirement folder"
  exit 0
fi

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

# shellcheck disable=SC2086
exec go test $folders -count=1 -parallel 32 "$@"
