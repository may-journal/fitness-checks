# 0002 Checks Run Only Where They Apply

## Why

I run every check in any repo with nothing to configure, and each one judges only what it fits. The checks stay fast on big trees, long branches, many docs and large diffs.

## Measurement

Checks that judge only the files they apply to and finish fast at any size

-

Checks

## Requirements

- 0002.1
  - Given Developer has every check program the repo builds
    - When each program declares the files it judges and a release is built
      - Then its patterns compile, it requires only known things, and it reaches the files it judges, in any folder, and no others
- 0002.2
  - Given Developer has a repo with nothing a check can judge and no checks chosen
    - When the installed checks run in that repo
      - Then each passes with zero files, every check is listed once in order, and `local-branch` declares itself and runs only inside fitness
- 0002.3
  - Given Developer has files, history or services that a check cannot read
    - When fitness runs that check
      - Then it skips them, names them as info, and still judges everything else
- 0002.4
  - Given Developer has a run that enables some checks
    - When fitness runs `read-repo-first`
      - Then it passes with zero files, shows the reminder and a table of the enabled checks, and keeps the table readable
- 0002.5
  - Given Developer has docs with and without fitness front matter
    - When fitness runs `no-fitness-front-matter`
      - Then it fails once for each file that holds a fitness key and passes the rest
- 0002.6
  - Given Developer has a large tree, a long branch, hundreds of docs or hundreds of changed files
    - When the tree, commit, prose and requirement checks judge it
      - Then each gives its verdict with only a few git calls in all
