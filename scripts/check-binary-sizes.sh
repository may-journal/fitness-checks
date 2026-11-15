#!/bin/bash
# cspell:ignore trimpath targz ldflags
# Fails when a check program built for a release grows more than the allowed percent past its size in tests/size-baseline.tsv.
# Usage: scripts/check-binary-sizes.sh, or with --update to write the baseline, which SIZE_BASELINE can move.
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"

baseline="${SIZE_BASELINE:-tests/size-baseline.tsv}"
allowed=10
work="$(mktemp -d)"
trap 'rm -rf "$work"' EXIT
# The baseline is for linux on arm64, built as a release is, so every machine measures the same program.
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -trimpath -buildvcs=false -ldflags="-s -w" -o "$work/" ./cmd/...

sizes="$work/sizes.tsv"
for program in "$work"/fitness-check-*; do
  printf '%s\t%s\n' "$(basename "$program")" "$(wc -c < "$program" | tr -d ' ')"
done | sort > "$sizes"

if [ "${1:-}" = --update ]; then
  { printf '# Sizes in bytes of each check program built for linux on arm64 with -trimpath and -s -w.\n'; cat "$sizes"; } > "$baseline"
  echo "check-binary-sizes: wrote $baseline"
  exit 0
fi

failed=0
while IFS=$'\t' read -r name bytes; do
  want="$(awk -F'\t' -v n="$name" '$1 == n { print $2 }' "$baseline")"
  if [ -z "$want" ]; then
    echo "SIZE: $name has no baseline, run scripts/check-binary-sizes.sh --update" >&2
    failed=1
  elif [ "$bytes" -gt $((want + want * allowed / 100)) ]; then
    echo "SIZE: $name is $bytes bytes, more than $allowed% over its baseline of $want" >&2
    failed=1
  fi
done < "$sizes"
echo "check-binary-sizes: $(wc -l < "$sizes" | tr -d ' ') programs checked against $baseline"
exit "$failed"
