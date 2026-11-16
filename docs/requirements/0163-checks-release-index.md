# 0163 Checks Release Index

## Why

An installer asks for one check at a major version, such as `fitness-check-gofmt-v1`, and needs the exact release and its files. A moving major tag per check and one index of every release answer that. The [README](../../README.md#release-index) records the index format.

## Measurement

Checks whose newest major tag and index entry name the same release
-
Checks

## Requirements

- 0163.1
    - Given Developer has checks released at 1.0.0, 1.2.0 and 2.0.0, and a major tag behind, moved through the API
        - When `scripts/release-major.sh` runs, then runs again
            - Then each `fitness-check-<name>-vN` tag, even a missing one, points at the newest release of major N, and the second run moves nothing
- 0163.2
    - Given Developer has a major tag `fitness-check-<name>-v1` on a later commit than its release 1.0.0
        - When `scripts/release-plan.sh` runs after a `fix` to that check
            - Then it plans 1.0.1, never reading the major tag as the last release
- 0163.3
    - Given Developer has a draft, a prerelease, a non-check release and two versions of one check
        - When `scripts/release-index.sh` reads them
            - Then the index holds the published checks with `latest`, `majors` and each artifact URL and sha256
- 0163.4
    - Given Developer has a release whose major tag is missing, and GitHub refuses a tag move as forbidden
        - When `scripts/release-major.sh` runs
            - Then it exits 1 with GitHub's message and creates no tag
- 0163.5
    - Given Developer has a major tag on a commit that is no ancestor of the release
        - When `scripts/release-major.sh` runs
            - Then the tag moves to the release, since the move is forced
- 0163.6
    - Given Developer has a release whose major tag exists, and GitHub refuses to move it
        - When `scripts/release-major.sh` runs
            - Then it exits 1 naming the tag, and never tries to create it
- 0163.7
    - Given Developer has two releases of a check, only the newer with a describe line that sets `files`, `exclude` and `requires`
        - When `scripts/release-index.sh` reads them
            - Then the newer version holds those fields and no empty or other describe field, and the older holds none
- 0163.8
    - Given Developer has a published index holding one version, and a release list with a newer one
        - When `scripts/release-index.sh` builds, with that index, then with an unreadable, empty or field-less one
            - Then it downloads only the new version, then both each time
