#!/bin/sh
# Prints, oldest first, each commit since a base that changed-checks.sh says can
# alter one check, with its subject: scripts/check-commits.sh <check> <base>.
set -eu
here="$(cd "$(dirname "$0")" && pwd)"
cd "$(git rev-parse --show-toplevel)"

check="${1:?usage: scripts/check-commits.sh <check name> <base>}"
base="${2:?usage: scripts/check-commits.sh <check name> <base>}"
"$here/changed-checks.sh" --each "$base" | while read -r commit checks; do
  case "$checks" in
  *\""$check"\"*) git log -1 --format='%H %s' "$commit" ;;
  esac
done
