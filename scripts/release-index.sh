#!/bin/sh
# cspell:ignore ltrimstr rtrimstr
# Prints index.json, which maps each check and major version to an exact release and its artifacts, from the published releases.
# Its format is in README.md, under Release index, and a file of releases JSON can stand in for GitHub as the first argument.
set -eu

releases="${1:-}"
if [ -z "$releases" ]; then
  releases="$(mktemp)"
  trap 'rm -f "$releases"' EXIT
  gh api --paginate --slurp "repos/$GITHUB_REPOSITORY/releases?per_page=100" | jq 'add' > "$releases"
fi

repository="${GITHUB_REPOSITORY:-may-journal/fitness-checks}"
jq --arg repository "$repository" '
  def number: split(".") | map(tonumber);
  def artifacts($tag): [.assets[] | select(.name | endswith(".tar"))
    | {key: (.name | ltrimstr($tag + "-") | rtrimstr(".tar")),
       value: {url: .browser_download_url, sha256: (.digest | ltrimstr("sha256:"))}}] | from_entries;
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
        versions: (map({key: .version, value: {tag, artifacts}}) | from_entries)}})
  | {format: 1, repository: $repository, checks: from_entries}' "$releases"
