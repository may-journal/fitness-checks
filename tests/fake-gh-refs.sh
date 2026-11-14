#!/bin/sh
# Stands in for gh in a test, applying its tag ref calls to the bare repo named by FAKE_REMOTE.
# It logs each call to calls.log there, and a PATCH of a missing tag fails as GitHub's does.
set -eu
printf '%s\n' "$*" >> "$FAKE_REMOTE/calls.log"
method=PATCH
sha=""
ref=""
for arg in "$@"; do
  case "$arg" in
  POST) method=POST ;;
  sha=*) sha="${arg#sha=}" ;;
  ref=*) ref="${arg#ref=}" ;;
  repos/*/git/refs/tags/*) ref="refs/tags/${arg##*/git/refs/tags/}" ;;
  esac
done
if [ "$method" = PATCH ]; then
  git -C "$FAKE_REMOTE" rev-parse -q --verify "$ref" > /dev/null || exit 1
fi
git -C "$FAKE_REMOTE" update-ref "$ref" "$sha"
