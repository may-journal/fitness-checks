#!/bin/sh
# Reads scripts/changed-checks.sh --each lines on stdin and prints, oldest first,
# each commit that can alter one check with its subject: check-commits.sh <check>.
set -eu

check="${1:?usage: scripts/check-commits.sh <check name> < changed-checks lines}"
while read -r commit checks; do
  case "$checks" in
  *\""$check"\"*) git log -1 --format='%H %s' "$commit" ;;
  esac
done
