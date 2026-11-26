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
   3. GIVEN it has unreadable services
      1. WHEN fitness runs that check
         1. THEN it skips them and names them as info
         2. THEN it judges the rest
4. Removed in 2026.11.26.0000: moved to the read-repo-first check's own requirement
5. Removed in 2026.11.26.0300: moved to the no-fitness-front-matter doc
6. Large repos stay fast
   1. GIVEN Developer has a large tree with a long branch
   2. GIVEN the tree has hundreds of docs and changed files
      1. WHEN the repo checks judge it
         1. THEN each gives its verdict with only a few git calls in all
7. Monorepos own packages
   1. GIVEN Developer has a monorepo whose `.fitnessrc.json` names its package folders
   2. GIVEN another repo names none
      1. WHEN the requirement checks judge its docs and tests
         1. THEN each package with code owns its docs and tests
         2. THEN 0001 repeats only across packages
         3. THEN a repo naming none is one
