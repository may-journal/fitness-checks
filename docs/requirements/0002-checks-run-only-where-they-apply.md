# 0002 Checks Run Only Where They Apply

## Why

I run every check in any repo with nothing to configure, and each one judges only what it fits. The checks stay fast on big trees, long branches, many docs and large diffs.

## Measurement

Checks that judge only the files they apply to and finish fast at any size

-

Checks

## Requirements

1. Programs declare their files
   1. GIVEN Developer has every check program the repo builds
   2. GIVEN a release is built
      1. WHEN each program declares the files it judges
         1. THEN its patterns compile and require only known things
         2. THEN it reaches the files it judges in any folder and no others
2. Empty repos pass
   1. GIVEN Developer has a repo a check cannot judge
   2. GIVEN no checks are chosen
      1. WHEN the installed checks run in that repo
         1. THEN each passes with zero files
         2. THEN every check is listed once
3. Unreadable inputs are skipped
   1. GIVEN Developer has files a check cannot read
   2. GIVEN it has unreadable history
      1. WHEN fitness runs that check
         1. THEN it skips them and names them as info
         2. THEN it judges the rest
4. Removed in 2026.11.26.0000: moved to the read-repo-first check's own requirement
5. Removed in 2026.11.26.0400: moved to the no-fitness-front-matter check's own requirement
6. Large repos stay fast
   1. GIVEN Developer has a large tree with a long branch
   2. GIVEN the tree has hundreds of docs and changed files
      1. WHEN the repo checks judge it
         1. THEN each gives its verdict with only a few git calls in all
7. Removed in 2026.11.26.0400: moved to the requirements and requirements-changed checks' own requirements
8. Bullet docs fail once
   1. GIVEN Developer has a requirement doc in the retired bullet form
      1. WHEN every check runs
         1. THEN exactly one check fails at the doc line
         2. THEN the others pass
