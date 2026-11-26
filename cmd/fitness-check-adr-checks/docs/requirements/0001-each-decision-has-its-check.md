# 0001 Each Decision Has Its Check

## Why

I trust each architecture decision because a check proves it and links its `ADR` issue.

## Measurement

ADR checks that link their issue and run

-

ADR checks

## Requirements

```gherkin
SCENARIO: 0001.1
  GIVEN Developer has a repo with or without `fitness/adr-NNNN` check folders
  WHEN the repo has no `fitness/` folder
  THEN it passes with no files
  WHEN each README links its `ADR` issue by full URL on a line of its own
  THEN it passes
    AND the paired check runs
  WHEN a README links no issue
    AND a README links one only inside a sentence
  THEN it fails
    BUT an `adrAsIssue` allow glob lets it wait
  WHEN an ADR folder has no check program
  THEN the run fails before any check starts
    BUT a superseded ADR is spared
  WHEN a paired check judges no files
  THEN it fails as a broken scan

SCENARIO: 0001.2
  GIVEN Developer has a commit or pull request changing an ADR check's files
  WHEN its README links another `ADR` issue
    AND no other file of its check changed
  THEN it fails naming the README
    AND it names the check
    AND it does so staged or in a pull request
  WHEN only the check's code changes
    AND the README's own text changes
  THEN it passes
  WHEN a README that linked no issue at the base links one
  THEN it passes without a code change
  WHEN a full scan runs with nothing staged
  THEN it judges no changes
```
