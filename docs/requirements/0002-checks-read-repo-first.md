# 0002 Checks Read Repo First

## Why

I see which checks guard this repo before I commit, so I read their rules instead of guessing.

## Measurement

runs that show the reminder and pass
-
runs of read-repo-first

## Requirements

- 0002.1
    - Given Developer has any repo
        - When fitness runs read-repo-first
            - Then it passes having judged 0 files
- 0002.2
    - Given Developer has any repo
        - When fitness runs read-repo-first
            - Then it asks whether I read the repo's decisions and enabled checks
- 0002.3
    - Given Developer has any repo
        - When fitness runs read-repo-first
            - Then it warns never to use `--no-verify`
- 0002.4
    - Given Developer has a run that enables only some checks
        - When read-repo-first prints its table
            - Then it lists exactly those checks, such as `duplicate-code`, beside their README paths
- 0002.5
    - Given Developer has a check name or path too long for its column
        - When read-repo-first prints its table
            - Then the cell is shortened, ending in `…`
- 0002.6
    - Given Developer has `NO_COLOR` is set
        - When fitness runs read-repo-first
            - Then its table prints with no color codes
