# 0164 Checks Parity

## Why

A check moved out of the engine must judge as the built-in did, so moving it changes no verdict. [scripts/check-parity.sh](../../scripts/check-parity.sh) runs both on fixed sample repos and compares the JSON.

## Measurement

Check and sample pairs where the binary and the built-in agree
-
Check and sample pairs

## Requirements

- 0164.1
    - Given Developer has every check built as a binary and as the engine's built-in
        - When `scripts/check-parity.sh` runs on the dirty sample
            - Then it exits 0, and `results.tsv` has a row for each built-in check, local-branch skipped, with findings from many checks
- 0164.2
    - Given Developer has a binary that reports a different verdict from its built-in
        - When `scripts/check-parity.sh` runs on the dirty sample
            - Then it exits 1 and names the check and the sample
- 0164.3
    - Given Developer has a built-in that disagrees with itself between two runs
        - When `scripts/check-parity.sh` runs
            - Then the pair is recorded as unstable, and fails only if the binary matches neither run on exit code and ok flag
- 0164.4
    - Given Developer has a built-in that disagrees with itself, and a binary whose ok flag matches neither run
        - When `scripts/check-parity.sh` runs
            - Then it exits 1 and names the check and the sample
