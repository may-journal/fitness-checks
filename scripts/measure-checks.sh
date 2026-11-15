#!/bin/bash
# cspell:ignore trimpath targz ldflags
# Records each check program's size before and after stripping, its archive sizes, and the timing of its run and start-up.
# Usage: scripts/measure-checks.sh [out-dir], with MEASURE_REPO for the repo to judge and MEASURE_RUNS for the suite runs, three by default.
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"

out="${1:-out/measure}"
case "$out" in /*) ;; *) out="$(pwd)/$out" ;; esac
case "$(pwd)/" in "$out/"*) echo "measure-checks: $out holds this checkout" >&2; exit 1 ;; esac
repo="${MEASURE_REPO:-$(pwd)}"
rm -rf "$out"
mkdir -p "$out/plain" "$out/stripped" "$out/bundle"

# ms prints the milliseconds a command takes to run, with its own output dropped.
ms() {
  perl -MTime::HiRes=time -e '$t = time(); open($o, ">&", \*STDOUT); open(STDOUT, ">", "/dev/null"); open(STDERR, ">", "/dev/null"); $ok = system(@ARGV) == 0; printf $o "%d\n", (time() - $t) * 1000; exit($ok ? 0 : 1)' -- "$@" < /dev/null
}

# The old bundle is the engine with every check built in, and the checks before were built unstripped.
go build -trimpath -buildvcs=false -ldflags="-s -w" -o "$out/bundle/" github.com/may-journal/fitness-runner/go/cmd/fitness
go build -trimpath -buildvcs=false -o "$out/plain/" ./cmd/...
go build -trimpath -buildvcs=false -ldflags="-s -w" -o "$out/stripped/" ./cmd/...
cp "$out/bundle/fitness" "$out/stripped/"

printf 'check\tplain_bytes\tstripped_bytes\ttar_bytes\ttargz_bytes\tdescribe_ms\n' > "$out/sizes.tsv"
for program in "$out"/stripped/fitness-check-*; do
  name="$(basename "$program")"
  describe="$(ms "$program" --describe)"
  tar -C "$out/stripped" -cf "$out/$name.tar" "$name"
  gzip -9 -n -k "$out/$name.tar"
  printf '%s\t%s\t%s\t%s\t%s\t%s\n' "${name#fitness-check-}" "$(wc -c < "$out/plain/$name" | tr -d ' ')" "$(wc -c < "$program" | tr -d ' ')" \
    "$(wc -c < "$out/$name.tar" | tr -d ' ')" "$(wc -c < "$out/$name.tar.gz" | tr -d ' ')" "$describe" >> "$out/sizes.tsv"
  rm -f "$out/$name.tar" "$out/$name.tar.gz"
done

# The suite runs each way, with the bundle alone and with the programs beside the engine, after one run that warms the tools.
(cd "$repo" && "$out/bundle/fitness" --all > /dev/null 2>&1 < /dev/null) 
printf 'run\tbundle_ms\tprograms_ms\n' > "$out/timing.tsv"
for run in $(seq "${MEASURE_RUNS:-3}"); do
  bundle="$(cd "$repo" && ms "$out/bundle/fitness" --all)"
  programs="$(cd "$repo" && ms "$out/stripped/fitness" --all)"
  printf '%s\t%s\t%s\n' "$run" "$bundle" "$programs" >> "$out/timing.tsv"
done

awk -F'\t' 'NR > 1 { p += $2; s += $3; t += $4; z += $5; d += $6; if ($6 > m) m = $6 }
  END { printf "checks %d: plain %d KiB, stripped %d KiB, tar %d KiB, tar.gz %d KiB, describe mean %.0f ms, max %d ms\n", NR - 1, p / 1024, s / 1024, t / 1024, z / 1024, d / (NR - 1), m }' "$out/sizes.tsv"
printf 'bundle: %d bytes\n' "$(wc -c < "$out/bundle/fitness" | tr -d ' ')"
awk -F'\t' 'NR > 1 { b += $2; p += $3 } END { printf "suite mean: bundle %.0f ms, programs %.0f ms\n", b / (NR - 1), p / (NR - 1) }' "$out/timing.tsv"
