#!/bin/sh
# Writes tests/requirement-times.tsv, the seconds the cases of each root
# requirement test took together, which scripts/requirement-shard.sh balances by.
# Usage: scripts/requirement-times.sh <log>..., each the output of
# tests/requirements.sh -v from one run, such as every requirements job of one CI run.
set -eu
cd "$(dirname "$0")/.."
[ "$#" -gt 0 ] || { echo "usage: scripts/requirement-times.sh <go test -v log>..." >&2; exit 1; }

# A parallel test reports no time of its own, so its cases one level down,
# such as Test0001_2/its_case, add up to its time. A CI log line opens with a timestamp.
times="$(sed -nE 's/^(.*Z )?    --- (PASS|FAIL): (Test[A-Za-z0-9_]+)\/[^ ]+ \(([0-9.]+)s\)$/\3	\4/p' "$@")"
[ -n "$times" ] || { echo "requirement-times.sh: the logs hold no case times" >&2; exit 1; }
{
  echo "# Seconds the cases of each root requirement test took together in CI; scripts/requirement-times.sh writes it."
  printf '%s\n' "$times" | awk -F'\t' '{ sum[$1] += $2 } END { for (t in sum) printf "%s\t%d\n", t, sum[t] + 0.5 }' | sort
} > tests/requirement-times.tsv
echo "requirement-times.sh: wrote $(grep -vc '^#' tests/requirement-times.tsv) test times to tests/requirement-times.tsv"
