#!/bin/sh
# Picks one share of the requirement suite, so CI runs the suite as parallel jobs.
# Usage: scripts/requirement-shard.sh <i>/<n> prints the go test -run pattern
# for share i of n of the root requirements package. Its top-level tests are
# dealt out in name order, so each runs in exactly one share.
# Usage: scripts/requirement-shard.sh checks <JSON array of check names>
# prints the requirement folders those checks own, one per line, or every
# check's folder when no array is given.
set -eu
cd "$(dirname "$0")/.."

if [ "${1:-}" = checks ]; then
  names="$(printf '%s' "${2:-}" | tr -d '[]"' | tr ',' ' ')"
  [ -n "${2:-}" ] || names="$(ls -d cmd/fitness-check-* | sed 's|^cmd/fitness-check-||')"
  for name in $names; do
    printf '%s\n' "$name" | grep -qxE '[a-z0-9-]+' || {
      echo "requirement-shard.sh: $name is not a check name" >&2
      exit 1
    }
    if [ -d "cmd/fitness-check-$name/requirements" ]; then
      echo "./cmd/fitness-check-$name/requirements"
    fi
  done
  exit 0
fi

share="${1:?usage: scripts/requirement-shard.sh <i>/<n> or checks <JSON array>}"
i="${share%/*}"
n="${share#*/}"
case "$i/$n" in
*[!0-9/]* | /* | */) echo "requirement-shard.sh: $share is not <i>/<n>" >&2; exit 1 ;;
esac
[ "$i" -ge 1 ] && [ "$i" -le "$n" ] || {
  echo "requirement-shard.sh: share $i is not between 1 and $n" >&2
  exit 1
}

# The test names come from the source, so no test binary is built to pick.
names="$(grep -ho '^func Test[A-Za-z0-9_]*' requirements/*_test.go | sed 's/^func //' | grep -vx TestMain | sort)"
printf '%s\n' "$names" | awk -v i="$i" -v n="$n" '
  (NR - 1) % n == i - 1 { picked = picked sep $0; sep = "|" }
  END { printf "^(%s)$\n", picked }'
