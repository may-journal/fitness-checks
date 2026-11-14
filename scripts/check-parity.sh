#!/bin/bash
# Runs each check as its binary and as the pinned engine's built-in on each sample, and fails when a verdict differs.
# Usage: scripts/check-parity.sh [out-dir], with CHECK_PARITY_SAMPLES and CHECK_PARITY_BIN to swap the sample list or the programs.
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"

out="${1:-out/parity}"
case "$out" in /*) ;; *) out="$(pwd)/$out" ;; esac
rm -rf "$out"
mkdir -p "$out/bin" "$out/samples"
if [ -n "${CHECK_PARITY_BIN:-}" ]; then
  cp "$CHECK_PARITY_BIN"/* "$out/bin/"
else
  go build -o "$out/bin/" github.com/may-journal/fitness-runner/go/cmd/fitness
  go build -o "$out/bin/" ./cmd/...
fi
fitness="$out/bin/fitness"

# A sample is a checkout of one pinned commit, as the engine meets a repo.
while read -r name url sha; do
  case "$name" in '' | '#'*) continue ;; esac
  dir="$out/samples/$name"
  git init -q "$dir"
  if [ "$url" = script ]; then
    "$sha" "$dir"
    git -C "$dir" add -A
    git -C "$dir" -c user.name=parity -c user.email=parity@example.test commit -q -m "added stuff"
    continue
  fi
  git -C "$dir" fetch -q --depth 1 "$url" "$sha"
  git -C "$dir" checkout -q FETCH_HEAD
done < "${CHECK_PARITY_SAMPLES:-tests/parity/samples.txt}"

# run prints a side's exit code on its first line and its normalized stdout after it.
run() {
  local code=0 stdout
  stdout="$("$@" 2> /dev/null < /dev/null)" || code=$?
  printf '%s\n%s\n' "$code" "$(printf '%s' "$stdout" | jq -S . 2> /dev/null || printf '%s' "$stdout")"
}

printf 'sample\tcheck\tstatus\tfindings\n' > "$out/results.tsv"
mismatches=0
for sample in "$out"/samples/*; do
  name="$(basename "$sample")"
  mkdir -p "$out/results/$name"
  for check in $("$fitness" check --list); do
    binary="$out/bin/fitness-check-$check"
    if [ ! -x "$binary" ]; then
      printf '%s\t%s\tskipped\t\n' "$name" "$check" >> "$out/results.tsv"
      continue
    fi
    built="$(cd "$sample" && run "$fitness" check "$check")"
    again="$(cd "$sample" && run "$fitness" check "$check")"
    bin="$(cd "$sample" && run "$binary" --root "$sample")"
    printf '%s\n' "$built" > "$out/results/$name/$check.builtin.json"
    printf '%s\n' "$bin" > "$out/results/$name/$check.binary.json"
    findings="$(printf '%s\n' "$built" | tail -n +2 | jq '.errors | length' 2> /dev/null || echo '?')"
    status=match
    if [ "$built" != "$again" ]; then
      # A built-in that disagrees with itself, as a test run can, has no verdict to match.
      status=unstable
    elif [ "$built" != "$bin" ]; then
      status=differ
      mismatches=$((mismatches + 1))
      echo "MISMATCH: check $check on sample $name, see $out/results/$name/$check.{builtin,binary}.json" >&2
    fi
    printf '%s\t%s\t%s\t%s\n' "$name" "$check" "$status" "$findings" >> "$out/results.tsv"
  done
done
count() { awk -F'\t' -v s="$1" '$3 == s' "$out/results.tsv" | wc -l | tr -d ' '; }
echo "check-parity: $(count match) matched, $mismatches differ, $(count unstable) unstable, $(count skipped) skipped; see $out/results.tsv"
[ "$mismatches" -eq 0 ]
