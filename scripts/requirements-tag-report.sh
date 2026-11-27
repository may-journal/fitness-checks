#!/bin/sh
# Writes the requirements failures a repo origin/main keeps after `--tag` runs, found in a throwaway clone.
# Usage: scripts/requirements-tag-report.sh <repo-dir> <output-file>
set -eu
here="$(cd "$(dirname "$0")/.." && pwd)"
repo="$(cd "${1:?usage: requirements-tag-report.sh <repo-dir> <output-file>}" && pwd)"
out="${2:?usage: requirements-tag-report.sh <repo-dir> <output-file>}"

bin="$(mktemp -d)"
tree="$(mktemp -d)"
cleanup() { rm -rf "$bin" "$tree"; }
trap cleanup EXIT

(cd "$here" && go build -o "$bin/requirements" ./cmd/fitness-check-requirements)
git -C "$repo" fetch -q origin main
git clone -q --shared --no-checkout "$repo" "$tree"
git -C "$tree" fetch -q "$repo" refs/remotes/origin/main
git -C "$tree" checkout -q --detach FETCH_HEAD

# The packages key of .fitnessrc.json names the package roots, one per line.
packages=""
if [ -f "$tree/.fitnessrc.json" ]; then
  packages="$(jq -r '.packages // [] | .[]' "$tree/.fitnessrc.json")"
fi
mkdir -p "$(dirname "$out")"
result="$(FITNESS_PACKAGES="$packages" "$bin/requirements" --root "$tree" --tag --dry-run)" || [ -n "$result" ] || { echo "the check crashed in $tree" >&2; exit 1; }
{
  echo "repo: $(basename "$repo") at $(git -C "$tree" rev-parse --short HEAD)"
  echo "failures left after --tag: $(printf '%s' "$result" | jq '.errors | length')"
  printf '%s' "$result" | jq -r '.errors[]'
} >"$out"
