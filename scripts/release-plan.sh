#!/bin/sh
# Prints as a JSON array each next tag, fitness-check-<name>-vX.Y.Z, for the
# checks HEAD alters since their last tag, or 1.0.0 for an untagged check.
set -eu
here="$(cd "$(dirname "$0")" && pwd)"
cd "$(git rev-parse --show-toplevel)"

# bump prints 3, 2, 1 or 0 for a major, minor, patch or no release. A "!" or a
# BREAKING CHANGE footer is major, feat minor, and fix, perf, revert or (deps) patch.
bump() {
  message="$(git log -1 --format=%B "$1")"
  subject="$(printf '%s\n' "$message" | head -n 1)"
  if printf '%s\n' "$subject" | grep -qE '^[a-z]+(\([^)]*\))?!:' ||
    printf '%s\n' "$message" | grep -qE '^BREAKING[ -]CHANGE:'; then
    echo 3
  elif printf '%s\n' "$subject" | grep -qE '^feat(\([^)]*\))?:'; then
    echo 2
  elif printf '%s\n' "$subject" | grep -qE '^((fix|perf|revert)(\([^)]*\))?|[a-z]+\(deps\)):|^Revert "'; then
    echo 1
  else
    echo 0
  fi
}

# Each distinct tagged commit's changed-checks.sh --each lines, read once each.
cache="$(mktemp -d)"
trap 'rm -rf "$cache"' EXIT
tags=""
for folder in cmd/fitness-check-*; do
  check="${folder#cmd/fitness-check-}"
  last="$(git tag -l "fitness-check-$check-v[0-9]*.[0-9]*.[0-9]*" --sort=-v:refname | head -n 1)"
  if [ -z "$last" ]; then
    tags="$tags fitness-check-$check-v1.0.0"
    continue
  fi
  base="$(git rev-parse "$last^{commit}")"
  [ -f "$cache/$base" ] || "$here/changed-checks.sh" --each "$base" > "$cache/$base"
  lines="$("$here/check-commits.sh" "$check" < "$cache/$base")"

  level=0
  for commit in $(printf '%s\n' "$lines" | cut -d' ' -f1); do
    this="$(bump "$commit")"
    [ "$this" -gt "$level" ] && level="$this"
  done
  version="${last##*-v}"
  major="${version%%.*}"
  rest="${version#*.}"
  minor="${rest%%.*}"
  patch="${rest#*.}"
  case "$level" in
  3) next="$((major + 1)).0.0" ;;
  2) next="$major.$((minor + 1)).0" ;;
  1) next="$major.$minor.$((patch + 1))" ;;
  *) continue ;;
  esac
  tags="$tags fitness-check-$check-v$next"
done

"$here/json-array.sh" $tags
