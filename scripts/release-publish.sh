#!/bin/sh
# Builds the check a release-plan.sh tag names, notes its commits and publishes
# it on HEAD, re-uploading on a rerun: scripts/release-publish.sh <tag>.
set -eu
here="$(cd "$(dirname "$0")" && pwd)"
cd "$(git rev-parse --show-toplevel)"

tag="${1:?usage: scripts/release-publish.sh <tag>}"
component="${tag%-v*}"
check="${component#fitness-check-}"
out="out/release/$tag"
"$here/release-build.sh" "$tag" "$out"

notes="$out.md"
last="$(git tag -l "$component-v*" --sort=-v:refname | grep -vxF "$tag" | head -n 1 || true)"
if [ -z "$last" ]; then
  echo "First release of $component on its own line." > "$notes"
else
  "$here/check-commits.sh" "$check" "$last" | sed 's/^/- /' > "$notes"
fi

if ! gh release view "$tag" > /dev/null 2>&1; then
  gh release create "$tag" --target "$(git rev-parse HEAD)" --title "$tag" --notes-file "$notes"
fi
gh release upload "$tag" "$out"/* --clobber
