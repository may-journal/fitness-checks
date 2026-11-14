#!/bin/sh
# Points each check's major tag, fitness-check-<name>-vN, at its highest release vN.x.y.
# Run it with no arguments in a checkout that holds every tag, and it pushes only the tags that move.
set -eu
cd "$(git rev-parse --show-toplevel)"

semver='^fitness-check-[a-z0-9-]+-v[0-9]+\.[0-9]+\.[0-9]+$'
seen=" "
for tag in $(git tag -l --sort=-v:refname | grep -E "$semver"); do
  component="${tag%-v*}"
  major="${tag##*-v}"
  major="${major%%.*}"
  move="$component-v$major"
  case "$seen" in *" $move "*) continue ;; esac
  seen="$seen$move "
  here="$(git rev-parse "$tag^{commit}")"
  [ "$(git rev-parse -q --verify "refs/tags/$move^{commit}" || true)" = "$here" ] && continue
  git tag -f "$move" "$here"
  git push --force origin "refs/tags/$move"
done
