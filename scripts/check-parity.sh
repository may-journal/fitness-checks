#!/bin/bash
# Runs each check as its binary and as the pinned engine's built-in on each sample, and fails when a verdict differs.
# Usage: scripts/check-parity.sh [out-dir], with CHECK_PARITY_SAMPLES and CHECK_PARITY_BIN to swap the sample list or the programs.
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"

out="${1:-out/parity}"
case "$out" in /*) ;; *) out="$(pwd)/$out" ;; esac
# Refuse a folder that holds this checkout, since the script clears it first.
case "$(pwd)/" in "$out/"*) echo "check-parity: $out holds this checkout" >&2; exit 1 ;; esac
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
while read -r name url sha || [ -n "$name" ]; do
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

# verdict prints the exit code and ok flag of a run's output.
verdict() { { printf '%s\n' "${1%%$'\n'*}"; printf '%s\n' "${1#*$'\n'}" | jq -r .ok 2> /dev/null || true; }; }

# run prints a side's exit code on its first line and its normalized stdout after it.
run() {
  local code=0 stdout
  stdout="$("$@" 2> /dev/null < /dev/null)" || code=$?
  printf '%s\n%s\n' "$code" "$(printf '%s' "${stdout//$PWD/<root>}" | jq -S . 2> /dev/null || printf '%s' "${stdout//$PWD/<root>}")"
}

# compare judges one check on one sample both ways, and writes its results.tsv row.
compare() {
  local name="$1" check="$2" sample="$out/samples/$1" status=match
  local binary="$out/bin/fitness-check-$check" built again bin findings
  if [ ! -x "$binary" ]; then
    printf '%s\t%s\tskipped\t\n' "$name" "$check" > "$out/rows/$name.$check"
    return 0
  fi
  built="$(cd "$sample" && run "$fitness" check "$check")"
  again="$(cd "$sample" && run "$fitness" check "$check")"
  bin="$(cd "$sample" && run "$binary" --root "$sample")"
  printf '%s\n' "$built" > "$out/results/$name/$check.builtin.json"
  printf '%s\n' "$bin" > "$out/results/$name/$check.binary.json"
  findings="$(printf '%s\n' "$built" | tail -n +2 | jq '.errors | length' 2> /dev/null || echo '?')"
  if [ "$built" != "$again" ]; then
    # A built-in that disagrees with itself, as a test run can, is held only to its exit code and ok flag.
    status=unstable
    [ "$(verdict "$bin")" = "$(verdict "$built")" ] || [ "$(verdict "$bin")" = "$(verdict "$again")" ] || status=differ
  elif [ "$built" != "$bin" ]; then
    status=differ
  fi
  [ "$status" != differ ] || echo "MISMATCH: check $check on sample $name, see $out/results/$name/$check.{builtin,binary}.json" >&2
  printf '%s\t%s\t%s\t%s\n' "$name" "$check" "$status" "$findings" > "$out/rows/$name.$check"
}
export -f compare run verdict
export out fitness

# Four pairs run at once, which keeps the wait short as checks and samples grow.
mkdir -p "$out/rows"
for sample in "$out"/samples/*; do
  mkdir -p "$out/results/$(basename "$sample")"
  for check in $("$fitness" check --list); do
    printf '%s %s\n' "$(basename "$sample")" "$check"
  done
done | xargs -P 4 -n 2 bash -c 'compare "$@"' _

printf 'sample\tcheck\tstatus\tfindings\n' > "$out/results.tsv"
cat "$out"/rows/* | sort >> "$out/results.tsv"
count() { awk -F'\t' -v s="$1" '$3 == s' "$out/results.tsv" | wc -l | tr -d ' '; }
echo "check-parity: $(count match) matched, $(count differ) differ, $(count unstable) unstable, $(count skipped) skipped; see $out/results.tsv"
[ "$(count differ)" -eq 0 ]
