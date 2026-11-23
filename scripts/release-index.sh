#!/bin/sh
# cspell:ignore ltrimstr rtrimstr slurpfile
# Prints index.json, the format README.md gives under Release index, from the published releases and their describe lines.
# Usage: scripts/release-index.sh [releases.json [describes.json [published.json]]], files or PUBLISHED_INDEX standing in for GitHub.
set -eu

work="$(mktemp -d)"
trap 'rm -rf "$work"' EXIT
repository="${GITHUB_REPOSITORY:-may-journal/fitness-checks}"
releases="${1:-}"
# Only a run without releases.json fetches the published index, so tests stay offline.
fetch=""
if [ -z "$releases" ]; then
  fetch=1
  releases="$work/releases.json"
  gh api --paginate --slurp "repos/$repository/releases?per_page=100" | jq 'add' > "$releases"
fi

# The one rule for which versions the index lists: the newest of each major, the only one a major pin selects.
# The index then grows with checks, not releases, and an older exact pin installs from its own tag.
kept='def number: split(".") | map(tonumber);
  def releases: [.[] | select((.draft | not) and (.prerelease | not)) | . as $release | .tag_name as $tag
    | ($tag | capture("^fitness-check-(?<name>[a-z0-9-]+)-v(?<version>[0-9]+\\.[0-9]+\\.[0-9]+)$")?)
    | select(. != null) | . + {tag: $tag, release: $release}];
  def kept: releases | group_by([.name, (.version | number | .[0])]) | map(max_by(.version | number));'

# The one rule for which describe fields count.
pick='def pick: {files, exclude, needs, requires} | with_entries(select(.value != null and .value != "" and .value != []));'

describes="${2:-}"
if [ -z "$describes" ]; then
  describes="$work/describes.json"
  # The published index already holds the fields of the versions it has seen.
  published="${3:-${PUBLISHED_INDEX:-}}"
  if [ -z "$published" ]; then
    published="$work/published.json"
    echo '{}' > "$published"
    if [ -n "$fetch" ]; then
      curl -fsSL --retry 3 -o "$published" "https://github.com/$repository/releases/download/fitness-check-index/index.json" 2> /dev/null || : > "$published"
    fi
  fi
  # A missing, empty or unreadable index means every version downloads; a version without fields does too.
  { jq -c "$pick"' [.checks[]?.versions[]? | {key: .tag, value: pick} | select(.value | length > 0)] | from_entries' "$published" > "$work/known.json" 2> /dev/null && [ -s "$work/known.json" ]; } || echo '{}' > "$work/known.json"
  jq -r --slurpfile known "$work/known.json" "$kept"' kept[] | .tag as $tag
    | select($known[0] | has($tag) | not)
    | .release.assets[]? | select(.name == $tag + ".describe.json") | "\($tag) \(.browser_download_url)"' "$releases" > "$work/assets"
  # A loop fed from a file, not a pipe, so a failed download stops the script.
  while read -r tag url; do
    curl -fsSL --retry 3 -o "$work/describe.json" "$url"
    jq -c --arg tag "$tag" '{($tag): .}' "$work/describe.json"
  done < "$work/assets" > "$work/parts"
  jq -s --slurpfile known "$work/known.json" '$known[0] + (add // {})' "$work/parts" > "$describes"
fi

jq --arg repository "$repository" --slurpfile described "$describes" "$kept$pick"'
  def fields($tag): ($described[0][$tag] // {}) | pick;
  def assets($tag; $suffix): [.assets[] | select(.name | endswith($suffix))
    | {key: (.name | ltrimstr($tag + "-") | rtrimstr($suffix)),
       value: {url: .browser_download_url, sha256: (.digest | ltrimstr("sha256:"))}}] | from_entries;
  def artifacts($tag): assets($tag; ".tar") as $tar | assets($tag; ".tar.gz") as $gz
    | $tar | with_entries(.value += (if $gz[.key] then {gz: $gz[.key]} else {} end));
  [kept[] | . + {artifacts: (.tag as $tag | .release | artifacts($tag))}]
  | group_by(.name)
  | map(sort_by(.version | number) | {
      key: .[0].name,
      value: {
        latest: .[-1].version,
        majors: (map({key: (.version | split(".")[0]), value: .version}) | from_entries),
        versions: (map({key: .version, value: ({tag, artifacts} + fields(.tag))}) | from_entries)}})
  | {format: 1, repository: $repository, checks: from_entries}' "$releases"
