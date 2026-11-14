# 0161 Checks Changed Checks

## Why

A change builds and tests only the checks it can alter, so one check's edit never waits on the other checks.

## Measurement

Changes whose CI runs exactly the checks they can alter
-
Changes

## Requirements

- 0161.1
    - Given Developer changes one file in one check's folder
        - When `scripts/changed-checks.sh` runs against the change's base
            - Then it names only that check
- 0161.2
    - Given Developer changes a check-only helper
        - When `scripts/changed-checks.sh` runs against the change's base
            - Then it names exactly the checks whose Go dependencies hold that helper
- 0161.3
    - Given Developer changes `go.mod`, or gives no base
        - When `scripts/changed-checks.sh` runs
            - Then it names every check
