#!/bin/sh
# Prints a run's summary for a job step: the job status and the count of each result in results.tsv.
# Usage: scripts/check-parity-summary.sh <results.tsv> <status>
set -eu

printf '### Job %s (%s)\n\n' "${GITHUB_JOB:-parity}" "$2"
[ -f "$1" ] || exit 0
awk -F'\t' 'NR > 1 { n[$3]++ } END { for (s in n) printf "- %s: %d\n", s, n[s] }' "$1" | sort
