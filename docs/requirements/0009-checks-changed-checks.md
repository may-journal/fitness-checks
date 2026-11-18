# 0009 Checks Changed Checks

## Why

A change builds and tests only the checks it can alter, so one check's edit never waits on the other checks.

## Measurement

Changes whose CI runs exactly the checks they can alter
-
Changes

## Requirements

- 0009.1
    - Given Developer changes one file in one check's folder
        - When `scripts/changed-checks.sh` runs against the change's base
            - Then it names only that check
- 0009.2
    - Given Developer changes a check-only helper, or a file it embeds
        - When `scripts/changed-checks.sh` runs against the change's base
            - Then it names exactly the checks whose Go dependencies hold that helper
- 0009.3
    - Given Developer changes `go.mod`, or gives no base
        - When `scripts/changed-checks.sh` runs
            - Then it names every check of the repo
- 0009.4
    - Given Developer has a check the change names
        - When `scripts/test-check.sh` runs with its name
            - Then it builds, vets and tests that check, and prints its `--describe` answer
