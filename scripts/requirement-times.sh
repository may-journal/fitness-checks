#!/bin/sh
# Writes tests/requirement-times.tsv, the seconds each root requirement test
# took, which scripts/requirement-shard.sh balances its shares by.
# Usage: scripts/requirement-times.sh <log>..., each the output of
# tests/requirements.sh -v, such as the logs of every requirements job in CI.
set -eu
cd "$(dirname "$0")/.."
[ "$#" -gt 0 ] || { echo "usage: scripts/requirement-times.sh <go test -v log>..." >&2; exit 1; }

times="$(sed -nE 's/^.*--- (PASS|FAIL): (Test[A-Za-z0-9_]+) \(([0-9.]+)s\).*$/\2	\3/p' "$@")"
[ -n "$times" ] || { echo "requirement-times.sh: the logs hold no top-level test times" >&2; exit 1; }
{
  echo "# Seconds each root requirement test took in CI; scripts/requirement-times.sh writes it."
  printf '%s\n' "$times" | awk -F'\t' '{ printf "%s\t%d\n", $1, $2 + 0.5 }' | sort -u -t"$(printf '\t')" -k1,1
} > tests/requirement-times.tsv
echo "requirement-times.sh: wrote $(grep -vc '^#' tests/requirement-times.tsv) test times to tests/requirement-times.tsv"
