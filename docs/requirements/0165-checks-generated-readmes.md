# 0165 Checks Generated READMEs

## Why

Each check's README says what its `doc.go` says, word for word, so the two never disagree.

## Measurement

Check READMEs equal to what `scripts/generate-readmes.sh` writes
-
Check READMEs

## Requirements

- 0165.1
    - Given Developer has the committed check READMEs and an empty Go module cache
        - When `scripts/generate-readmes.sh -check` runs
            - Then it passes, since each README is what the script writes
- 0165.2
    - Given Developer edits a check's doc comment and leaves its README alone
        - When `scripts/generate-readmes.sh -check` runs
            - Then it fails naming that README
- 0165.3
    - Given Developer edits a check's doc comment
        - When `scripts/generate-readmes.sh` runs
            - Then that README holds the new text under its title, and a second `-check` passes
