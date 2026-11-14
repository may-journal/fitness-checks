# 0162 Checks Release Per Check

## Why

Each check releases on its own line, so a consumer takes one check's fix without the others. `go.mod` and `go.sum` hold every version, and [scripts/changed-checks.sh](../../scripts/changed-checks.sh) decides which checks a change can alter, for CI and for release alike.

## Measurement

Pushes to main that release exactly the checks they can alter
-
Pushes to main

## Requirements

- 0162.1
    - Given Developer merges a `fix` commit to one file in one check's folder, after every check has a release tag
        - When `scripts/release-plan.sh` runs on main
            - Then it plans a patch release of only that check
- 0162.2
    - Given Developer merges a `feat` commit to a check-only helper
        - When `scripts/release-plan.sh` runs on main
            - Then it plans a minor release of exactly the checks whose Go dependencies hold that helper
- 0162.3
    - Given Developer merges a `(deps)` commit to `go.mod`
        - When `scripts/release-plan.sh` runs on main
            - Then it plans a patch release of every check
- 0162.4
    - Given Developer merges to a check a `docs` commit, a commit marked breaking with `!`, or a GitHub `Revert` commit
        - When `scripts/release-plan.sh` runs on main
            - Then it plans no release, a major release, and a patch release in turn
- 0162.5
    - Given Developer has a check with no release tag yet
        - When `scripts/release-plan.sh` runs on main
            - Then it plans that check's first release, 1.0.0, from the checks the repo holds
- 0162.6
    - Given Developer builds a release for one check
        - When `scripts/release-build.sh` runs for its tag
            - Then the output holds an archive and a checksum for linux and darwin on amd64 and arm64, named by tag
