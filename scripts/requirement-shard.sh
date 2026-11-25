#!/bin/sh
# Picks one share of the requirement suite, so CI runs the suite as parallel jobs.
# Usage: scripts/requirement-shard.sh <i>/<n> prints the go test -run pattern
# for share i of n of the root requirements package. Its top-level tests are
# dealt out by their recorded times, so each runs in exactly one share.
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
# Each test weighs the seconds tests/requirement-times.tsv gives it, or 10
# when it has no row. The heaviest goes first to the lightest share, so the
# shares end close together.
printf '%s\n' "$names" | awk -v i="$i" -v n="$n" -F'\t' '
  FNR == NR { if ($0 !~ /^#/) weight[$1] = $2; next }
  { w = ($1 in weight) ? weight[$1] : 10; printf "%012.3f\t%s\n", w, $1 }
' tests/requirement-times.tsv - | sort -t"$(printf '\t')" -k1,1r -k2,2 | awk -v i="$i" -v n="$n" -F'\t' '
  {
    lightest = 1
    for (s = 2; s <= n; s++) if (load[s] < load[lightest]) lightest = s
    load[lightest] += $1
    if (lightest == i) { picked = picked sep $2; sep = "|" }
  }
  END { printf "^(%s)$\n", picked }'
