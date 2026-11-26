# 0006 Changes Come Paired

## Why

I change a requirement test and its doc together, so neither drifts from the other.

## Measurement

changed requirement tests and docs that come with their partner

-

changed requirement tests and docs

## Requirements

1. A test and its doc change together
   1. GIVEN Developer has requirement tests with docs
      1. WHEN a test or a doc changes without its partner
         1. THEN it fails naming both files
      2. WHEN a test is named for two requirements
         1. THEN a change to either doc pairs it
      3. WHEN a shared helper changes
         1. THEN it needs no partner
      4. WHEN a merge or a cherry-pick or a revert replays a change
         1. THEN only its own changes are judged
      5. WHEN a commit is amended
         1. THEN it is judged since its parent
      6. WHEN a failing test is fixed
         1. THEN it needs no doc
      7. WHEN only comments or formatting change
         1. THEN it needs no partner
         2. THEN a reformat or rename of numbered acceptances needs no partner
