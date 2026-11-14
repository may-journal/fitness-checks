#!/bin/sh
# Uploads index.json to the fixed release fitness-check-index, which an installer reads at a stable URL.
# The release is created on the first run and keeps its one asset current after that.
set -eu
here="$(cd "$(dirname "$0")" && pwd)"
cd "$(git rev-parse --show-toplevel)"

out="out/index"
rm -rf "$out"
mkdir -p "$out"
"$here/release-index.sh" > "$out/index.json"
tag="fitness-check-index"
if gh release view "$tag" > /dev/null 2>&1; then
  gh release upload "$tag" "$out/index.json" --clobber
else
  gh release create "$tag" "$out/index.json" --target "$(git rev-parse HEAD)" --latest=false \
    --title "Release index" --notes "index.json resolves each check and major version to an exact release and its artifacts."
fi
