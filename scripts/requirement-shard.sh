#!/bin/sh
# Picks one share of the requirement suite so CI runs it as parallel jobs.
# Usage: scripts/requirement-shard.sh <i>/<n>, checks [JSON array of check names [<i>/<n>]], shares [JSON array], or suite [base].
set -eu

# suite prints true when a change since base touches the shared suite or the
# slow release scripts. It also does with no base.
if [ "${1:-}" = suite ]; then
  base="${2:-}"
  if [ -z "$base" ] || ! git cat-file -e "$base^{commit}" 2>/dev/null; then
    echo true
  elif git diff --name-only "$base" HEAD |
    grep -qE '^(go\.(mod|sum)|tests/|\.github/workflows/checks\.yml|scripts/(release-build|measure-checks|check-binary-sizes|requirement-shard)\.sh)'; then
    echo true
  else
    echo false
  fi
  exit 0
fi
cd "$(dirname "$0")/.."

# folders prints the requirement folders the named checks own, or every
# check's folder when no array is given.
folders() {
  names="$(printf '%s' "${1:-}" | tr -d '[]"' | tr ',' ' ')"
  [ -n "${1:-}" ] || names="$(ls -d cmd/fitness-check-* | sed 's|^cmd/fitness-check-||')"
  for name in $names; do
    printf '%s\n' "$name" | grep -qxE '[a-z0-9-]+' || {
      echo "requirement-shard.sh: $name is not a check name" >&2
      exit 1
    }
    if [ -d "cmd/fitness-check-$name/requirements" ]; then
      echo "./cmd/fitness-check-$name/requirements"
    fi
  done
}

# shares prints the shares, as a JSON array of <i>/<n>, that the folders of the
# named checks (every check with no array) are dealt into: one per folder, up
# to six, so a batch of changes runs in parallel jobs. It prints [] for none.
if [ "${1:-}" = shares ]; then
  all="$(folders "${2:-}")" || exit 1
  count="$(printf '%s\n' "$all" | grep -c . || true)"
  [ "$count" -le 6 ] || count=6
  labels=""
  i=1
  while [ "$i" -le "$count" ]; do
    labels="$labels $i/$count"
    i=$((i + 1))
  done
  # shellcheck disable=SC2086
  scripts/json-array.sh $labels
  exit 0
fi

# checks prints the folders of the named checks, or only share <i>/<n> of them. The heaviest folder, by its test files' size, goes first to the lightest share.
if [ "${1:-}" = checks ]; then
  if [ -z "${3:-}" ]; then
    folders "${2:-}"
    exit 0
  fi
  i="${3%/*}"
  n="${3#*/}"
  [ "$3" = "$i/$n" ] || { echo "requirement-shard.sh: $3 is not <i>/<n>" >&2; exit 1; }
  case "$3" in
  */*/* | */ | /* | *[!0-9/]*) echo "requirement-shard.sh: $3 is not <i>/<n>" >&2; exit 1 ;;
  esac
  [ "$i" -ge 1 ] && [ "$i" -le "$n" ] || {
    echo "requirement-shard.sh: share $i is not between 1 and $n" >&2
    exit 1
  }
  all="$(folders "${2:-}")" || exit 1
  for folder in $all; do
    printf '%012d\t%s\n' "$(cat "$folder"/*_test.go | wc -c)" "$folder"
  done | sort -t"$(printf '\t')" -k1,1r -k2,2 | awk -v i="$i" -v n="$n" -F'\t' '
    {
      lightest = 1
      for (s = 2; s <= n; s++) if (load[s] < load[lightest]) lightest = s
      load[lightest] += $1
      if (lightest == i) print $2
    }'
  exit 0
fi

# <i>/<n> prints the -run pattern for share i of n of the root requirements package.
share="${1:?usage: scripts/requirement-shard.sh <i>/<n> or checks <JSON array>}"
case "$share" in */*/* | */ | /*) echo "requirement-shard.sh: $share is not <i>/<n>" >&2; exit 1 ;; */*) ;; *) echo "requirement-shard.sh: $share is not <i>/<n>" >&2; exit 1 ;; esac
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
[ -f tests/requirement-times.tsv ] || { echo "requirement-shard.sh: tests/requirement-times.tsv is missing" >&2; exit 1; }
names="$(grep -ho '^func Test[A-Za-z0-9_]*' requirements/*_test.go | sed 's/^func //' | grep -vx TestMain | sort)"
# Each test weighs the seconds tests/requirement-times.tsv gives it, or their
# mean when it has no row. The heaviest goes first to the lightest share.
printf '%s\n' "$names" | awk -F'\t' '
  FNR == NR { if ($0 !~ /^#/) { weight[$1] = $2; sum += $2; rows++ }; next }
  { w = ($1 in weight) ? weight[$1] : (rows ? sum / rows : 1); printf "%012.3f\t%s\n", w, $1 }
' tests/requirement-times.tsv - | sort -t"$(printf '\t')" -k1,1r -k2,2 | awk -v i="$i" -v n="$n" -F'\t' '
  {
    lightest = 1
    for (s = 2; s <= n; s++) if (load[s] < load[lightest]) lightest = s
    load[lightest] += $1
    if (lightest == i) { picked = picked sep $2; sep = "|" }
  }
  END {
    if (picked == "") { print "requirement-shard.sh: share " i " of " n " holds no test" > "/dev/stderr"; exit 1 }
    printf "^(%s)$\n", picked
  }'
