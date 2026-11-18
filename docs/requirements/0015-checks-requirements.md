# 0015 Checks Requirements

## Why

Every requirement keeps its IDs and changes with its test, so a doc and its proof never drift apart.

## Measurement

changed requirement tests and docs whose partner changed too

-

changed requirement tests and docs

## Requirements

- 0015.8
  - Given Developer has an ID that an earlier commit deleted
    - When a doc brings it back
      - Then it fails
- 0015.9
  - Given Developer has a requirement test or doc changed without its partner
    - When fitness judges my commit or pull request
      - Then it fails naming both files and the missing half
- 0015.10
  - Given Developer has a changed test file named for two requirements
    - When fitness runs on staged changes
      - Then it fails naming both docs until either doc changes
- 0015.11
  - Given Developer has a changed shared test helper that proves no requirement
    - When fitness runs on staged changes
      - Then the check passes
- 0015.14
  - Given Developer has a merge, cherry-pick, revert, or base merge that replays an unpaired change
    - When fitness judges my commit or pull request
      - Then only the branch's own changes are judged
- 0015.15
  - Given Developer has a commit with a requirement test alone and its doc staged
    - When I amend the commit with the hooks
      - Then the commit goes through, judged on everything since its parent
- 0015.16
  - Given Developer has a requirement test that fails the check at its base
    - When I fix the test without changing its doc
      - Then my commit, amend, or pull request passes, while a test that passed at its base still fails without its doc
- 0015.17
  - Given Developer has a test changed only in comments, or a doc changed only in formatting, comments, name, or a deleted Source line
    - When fitness judges my commit or pull request
      - Then it needs no partner and counts as none, while editing a Source line still needs one
