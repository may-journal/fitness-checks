#!/bin/bash
# cspell:ignore trimpath targz ldflags
# Fails when a check program built for a release is more than the allowed percent off its size in tests/size-baseline.tsv.
# Usage: scripts/check-binary-sizes.sh, or with --update to write the baseline, and SIZE_BASELINE and SIZE_BUILD_DIR move the files.
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"

baseline="${SIZE_BASELINE:-tests/size-baseline.tsv}"
flags='-trimpath -buildvcs=false -ldflags="-s -w"'
allowed=10
# The baseline is for the flags the release build uses, so a changed release build must change this script too.
grep -qF -- "$flags" scripts/release-build.sh || { echo "SIZE: scripts/release-build.sh no longer builds with $flags" >&2; exit 1; }

work="${SIZE_BUILD_DIR:-}"
if [ -z "$work" ]; then
  work="$(mktemp -d)"
  trap 'rm -rf "$work"' EXIT
fi
mkdir -p "$work"
# The baseline is for linux on arm64 and the go.mod toolchain, so every machine measures the same program.
gover="$(awk '/^go / { print $2 }' go.mod)"
case "$gover" in *.*.*) ;; *) gover="$gover.0" ;; esac
export GOTOOLCHAIN="go$gover"
if ! ls "$work"/fitness-check-* > /dev/null 2>&1; then
  CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -trimpath -buildvcs=false -ldflags="-s -w" -o "$work/" ./cmd/...
fi

sizes="$work/sizes.tsv"
for program in "$work"/fitness-check-*; do
  printf '%s\t%s\n' "$(basename "$program")" "$(wc -c < "$program" | tr -d ' ')"
done | sort > "$sizes"

if [ "${1:-}" = --update ]; then
  { printf '# Sizes in bytes of each check program built for linux on arm64 with path trimming and stripping.\n'; cat "$sizes"; } > "$baseline"
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
  elif [ "$bytes" -lt $((want - want * allowed / 100)) ]; then
    echo "SIZE: $name is $bytes bytes, more than $allowed% under its baseline of $want, run scripts/check-binary-sizes.sh --update" >&2
    failed=1
  fi
done < "$sizes"
# A baseline row for a program no longer built is stale.
for name in $(grep -v '^#' "$baseline" | cut -f1); do
  [ -e "$work/$name" ] || { echo "SIZE: $name is in the baseline but is no longer built, run scripts/check-binary-sizes.sh --update" >&2; failed=1; }
done
echo "check-binary-sizes: $(wc -l < "$sizes" | tr -d ' ') programs checked against $baseline"
exit "$failed"
