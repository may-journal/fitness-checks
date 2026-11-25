# 0002 Checks Run Only Where They Apply

## Why

I run every check in any repo with nothing to configure, and each one judges only what it fits. The checks stay fast on big trees, long branches, many docs and large diffs.

## Measurement

Checks that judge only the files they apply to and finish fast at any size

-

Checks

## Requirements

```gherkin
SCENARIO: 0002.1
  GIVEN Developer has every check program the repo builds
  WHEN each program declares the files it judges and a release is built
  THEN its patterns compile
    AND it requires only known things
    AND it reaches the files it judges in any folder and no others

SCENARIO: 0002.2
  GIVEN Developer has a repo with nothing a check can judge and no checks chosen
  WHEN the installed checks run in that repo
  THEN each passes with zero files
    AND every check is listed once in order
    AND `local-branch` declares itself and runs only inside fitness

SCENARIO: 0002.3
  GIVEN Developer has files or history or services that a check cannot read
  WHEN fitness runs that check
  THEN it skips them
    AND it names them as info
    AND it still judges everything else

SCENARIO: 0002.4
  GIVEN Developer has a run that enables some checks
  WHEN fitness runs `read-repo-first`
  THEN it passes with zero files
    AND it shows the reminder and a table of the enabled checks
    AND it keeps the table readable

SCENARIO: 0002.5
  GIVEN Developer has docs with and without fitness front matter
  WHEN fitness runs `no-fitness-front-matter`
  THEN it fails once for each file that holds a fitness key and passes the rest

SCENARIO: 0002.6
  GIVEN Developer has a large tree or a long branch or hundreds of docs or hundreds of changed files
  WHEN the tree and commit and prose and requirement checks judge it
  THEN each gives its verdict with only a few git calls in all

SCENARIO: 0002.7
  GIVEN Developer has a monorepo whose `.fitnessrc.json` names its package folders or a repo that names none
  WHEN the requirement checks judge its docs and tests
  THEN each package with code owns its docs and tests and numbers
    AND 0001 repeats only across packages
    AND a repo naming none is one
```
