# 0001 Each Decision Has Its Check

## Why

I trust each architecture decision because a check proves it and links its `ADR` issue.

## Measurement

ADR checks that link their issue and run

-

ADR checks

## Requirements

- 0001.1
  - Given Developer has a repo, with or without `fitness/adr-NNNN` check folders
    1. When the repo has no `fitness/` folder
       - Then it passes with no files
    2. When each README links its `ADR` issue by full URL on a line of its own
       - Then it passes, and the paired check runs
    3. When a README links no issue, or one only inside a sentence
       - Then it fails, unless an `adrAsIssue` allow glob lets it wait
    4. When an ADR folder has no check program
       - Then the run fails before any check starts, unless the ADR is superseded
    5. When a paired check judges no files
       - Then it fails as a broken scan
- 0001.2
  - Given Developer has a commit or pull request changing an ADR check's files
    1. When its README links another `ADR` issue and no other file of its check changed
       - Then it fails naming the README and the check, staged or in a pull request
    2. When only the check's code, or the README's own text, changes
       - Then it passes
    3. When a README that linked no issue at the base links one
       - Then it passes without a code change
    4. When a full scan runs with nothing staged
       - Then it judges no changes
