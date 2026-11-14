#!/bin/sh
# Builds the check a release-plan.sh tag names and publishes it on HEAD with
# notes, finishing a draft a failed run left: release-publish.sh <tag>.
set -eu
here="$(cd "$(dirname "$0")" && pwd)"
cd "$(git rev-parse --show-toplevel)"

tag="${1:?usage: scripts/release-publish.sh <tag>}"
component="${tag%-v*}"
check="${component#fitness-check-}"
out="out/release/$tag"
"$here/release-build.sh" "$tag" "$out"

notes="$out.md"
last="$(git tag -l "$component-v[0-9]*" --sort=-v:refname | grep -vxF "$tag" | head -n 1 || true)"
if [ -z "$last" ]; then
  echo "First release of $component on its own line." > "$notes"
else
  each="$("$here/changed-checks.sh" --each "$last")"
  lines="$(printf '%s\n' "$each" | "$here/check-commits.sh" "$check")"
  printf '%s\n' "$lines" | sed 's/^/- /' > "$notes"
fi

# gh uploads the programs to a draft and publishes it last, hiding it until then.
# A rerun finishes the draft a failed upload left.
if gh release view "$tag" > /dev/null 2>&1; then
  gh release upload "$tag" "$out"/* --clobber
  gh release edit "$tag" --draft=false
else
  gh release create "$tag" "$out"/* --target "$(git rev-parse HEAD)" --title "$tag" --notes-file "$notes"
fi
