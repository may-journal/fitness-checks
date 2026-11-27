#!/bin/sh
# Diffs each check's --describe answer between its last release before a commit and the current tree.
# Usage: scripts/describe-diff.sh <batch commit> <check name>...
set -eu
cd "$(dirname "$0")/.."

batch="${1:?usage: scripts/describe-diff.sh <batch commit> <check name>...}"
shift
[ "$#" -gt 0 ] || { echo "describe-diff.sh: name at least one check" >&2; exit 2; }

work="$(mktemp -d "${TMPDIR:-/tmp}/describe-diff.XXXXXX")"
trap 'rm -rf "$work"' EXIT
differs=0

for name in "$@"; do
  # The newest release tag of the check that the batch's parent holds.
  tag="$(git tag --merged "$batch^" -l "fitness-check-$name-v[0-9]*.[0-9]*.[0-9]*" --sort=-v:refname | head -1)"
  [ -n "$tag" ] || { echo "describe-diff.sh: no release tag of $name before $batch" >&2; exit 2; }
  # A tag's tree builds in its own folder, so the build sees that release's code and go.mod.
  mkdir "$work/src-$name"
  git archive "$tag" | tar -x -C "$work/src-$name"
  (cd "$work/src-$name" && go build -o "$work/old-$name" "./cmd/fitness-check-$name")
  go build -o "$work/new-$name" "./cmd/fitness-check-$name"
  "$work/old-$name" --describe >"$work/old-$name.json"
  "$work/new-$name" --describe >"$work/new-$name.json"
  if diff -u --label "$name at $tag" --label "$name now" "$work/old-$name.json" "$work/new-$name.json"; then
    echo "$name: identical to $tag"
  else
    differs=1
  fi
done
exit "$differs"
