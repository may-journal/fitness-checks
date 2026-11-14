#!/bin/sh
# Stands in for gh in a test, applying its tag ref calls to the bare repo named by FAKE_REMOTE.
# It logs each call to calls.log there, and a PATCH of a missing tag fails as GitHub's does.
set -eu
printf '%s\n' "$*" >> "$FAKE_REMOTE/calls.log"
method=PATCH
sha=""
ref=""
force=""
for arg in "$@"; do
  case "$arg" in
  POST) method=POST ;;
  force=*) force="${arg#force=}" ;;
  sha=*) sha="${arg#sha=}" ;;
  ref=*) ref="${arg#ref=}" ;;
  repos/*/git/refs/tags/*) ref="refs/tags/${arg##*/git/refs/tags/}" ;;
  esac
done
if [ "$method" = PATCH ] && [ -n "${FAKE_PATCH_ERROR:-}" ]; then
  echo "$FAKE_PATCH_ERROR" >&2
  exit 1
fi
if [ "$method" = PATCH ]; then
  git -C "$FAKE_REMOTE" rev-parse -q --verify "$ref" > /dev/null || { echo "gh: Reference does not exist (HTTP 422)" >&2; exit 1; }
fi
git -C "$FAKE_REMOTE" cat-file -e "$sha^{commit}" 2> /dev/null || { echo "gh: Object does not exist (HTTP 422)" >&2; exit 1; }
if [ "$method" = PATCH ] && [ "$force" != true ]; then
  git -C "$FAKE_REMOTE" merge-base --is-ancestor "$ref" "$sha" || { echo "gh: Update is not a fast forward (HTTP 422)" >&2; exit 1; }
fi
git -C "$FAKE_REMOTE" update-ref "$ref" "$sha"
