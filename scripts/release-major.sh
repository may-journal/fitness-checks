#!/bin/sh
# Points each check's major tag, fitness-check-<name>-vN, at its highest release vN.x.y.
# Run it with no arguments in a checkout that holds every tag, with GH_TOKEN set, and it moves only the tags that differ.
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
  # The API moves a tag where a push from a workflow is refused for a commit that changes a workflow file.
  gh api -X PATCH "repos/$GITHUB_REPOSITORY/git/refs/tags/$move" -f sha="$here" -F force=true --silent 2> /dev/null ||
    gh api -X POST "repos/$GITHUB_REPOSITORY/git/refs" -f ref="refs/tags/$move" -f sha="$here" --silent
done
