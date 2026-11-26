# 0001 Each Decision Has Its Check

## Why

I trust each architecture decision because a check proves it and links its `ADR` issue.

## Measurement

ADR checks that link their issue and run

-

ADR checks

## Requirements

1. Checks run only for linked decisions
   1. GIVEN Developer has a repo
   2. GIVEN it may hold `fitness/adr-NNNN` check folders
      1. WHEN the repo has no `fitness/` folder
         1. THEN it passes with no files
      2. WHEN each README links its `ADR` issue by full URL on a line of its own
         1. THEN it passes
         2. THEN the paired check runs
      3. WHEN a README links no issue
         1. THEN it fails
         2. THEN a link only inside a sentence fails too
         3. THEN an `adrAsIssue` allow glob lets it wait
      4. WHEN an ADR folder has no check program
         1. THEN the run fails before any check starts
         2. THEN a superseded ADR is exempt
      5. WHEN a paired check judges no files
         1. THEN it fails as a broken scan
2. Changed links need changed checks
   1. GIVEN Developer has a commit or pull request changing an ADR check's files
      1. WHEN its README links another `ADR` issue and no other file of its check changed
         1. THEN it fails naming the README and the check
         2. THEN it fails the same way staged or in a pull request
      2. WHEN only the check's code changes
         1. THEN it passes
         2. THEN a change to only the README's own text passes too
      3. WHEN a README that linked no issue at the base links one
         1. THEN it passes without a code change
      4. WHEN a full scan runs with nothing staged
         1. THEN it judges no changes
