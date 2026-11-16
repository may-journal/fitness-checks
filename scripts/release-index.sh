#!/bin/sh
# cspell:ignore ltrimstr rtrimstr slurpfile
# Prints index.json, the format README.md gives under Release index, from the published releases and their describe lines.
# Usage: scripts/release-index.sh [releases.json [describes.json]], files that stand in for GitHub and the describe assets.
set -eu

work="$(mktemp -d)"
trap 'rm -rf "$work"' EXIT
releases="${1:-}"
if [ -z "$releases" ]; then
  releases="$work/releases.json"
  gh api --paginate --slurp "repos/$GITHUB_REPOSITORY/releases?per_page=100" | jq 'add' > "$releases"
fi

describes="${2:-}"
if [ -z "$describes" ]; then
  describes="$work/describes.json"
  jq -r '.[] | select((.draft | not) and (.prerelease | not)) | .tag_name as $tag
    | .assets[]? | select(.name == $tag + ".describe.json") | "\($tag) \(.browser_download_url)"' "$releases" > "$work/assets"
  # A loop fed from a file, not a pipe, so a failed download stops the script.
  while read -r tag url; do
    curl -fsSL --retry 3 -o "$work/describe.json" "$url"
    jq -c --arg tag "$tag" '{($tag): .}' "$work/describe.json"
  done < "$work/assets" > "$work/parts"
  jq -s 'add // {}' "$work/parts" > "$describes"
fi

repository="${GITHUB_REPOSITORY:-may-journal/fitness-checks}"
jq --arg repository "$repository" --slurpfile described "$describes" '
  def fields($tag): ($described[0][$tag] // {})
    | {files, exclude, needs, requires} | with_entries(select(.value != null and .value != "" and .value != []));
  def number: split(".") | map(tonumber);
  def assets($tag; $suffix): [.assets[] | select(.name | endswith($suffix))
    | {key: (.name | ltrimstr($tag + "-") | rtrimstr($suffix)),
       value: {url: .browser_download_url, sha256: (.digest | ltrimstr("sha256:"))}}] | from_entries;
  def artifacts($tag): assets($tag; ".tar") as $tar | assets($tag; ".tar.gz") as $gz
    | $tar | with_entries(.value += (if $gz[.key] then {gz: $gz[.key]} else {} end));
  [.[] | select((.draft | not) and (.prerelease | not)) | . as $release | .tag_name as $tag
    | ($tag | capture("^fitness-check-(?<name>[a-z0-9-]+)-v(?<version>[0-9]+\\.[0-9]+\\.[0-9]+)$")?)
    | select(. != null)
    | . + {tag: $tag, artifacts: ($release | artifacts($tag))}]
  | group_by(.name)
  | map(sort_by(.version | number) | {
      key: .[0].name,
      value: {
        latest: .[-1].version,
        majors: (map({key: (.version | split(".")[0]), value: .version}) | from_entries),
        versions: (map({key: .version, value: ({tag, artifacts} + fields(.tag))}) | from_entries)}})
  | {format: 1, repository: $repository, checks: from_entries}' "$releases"
