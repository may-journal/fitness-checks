# 0004 Code Kinds And Runs

## Why

I get each kind of code and each run judged alike, so a pair is asked for only where one is due.

## Measurement

product code changes judged by their kind and run

-

product code changes

## Requirements

1. Code is judged by its kind
   1. GIVEN Developer has a repo with requirement docs
   2. GIVEN it has Swift sources plus shell scripts
      1. WHEN Developer edits a Swift source or a shell script alone
         1. THEN it fails naming the file
      2. WHEN Developer pairs a Swift source with a Swift test plus a doc
         1. THEN it passes
      3. WHEN a directive or a shebang or a cgo preamble changes among comment edits
         1. THEN it counts as code
2. Each run judges its own changes
   1. GIVEN Developer has a repo with requirement docs
   2. GIVEN product code changes with no test or doc
      1. WHEN a pre-commit hook runs alone
         1. THEN it fails naming the file
      2. WHEN a commit-msg hook sits beside it
         1. THEN the pre-commit hook leaves the check out
         2. THEN the commit-msg hook fails naming the file
      3. WHEN an amend stages code whose parent changed the test plus doc
         1. THEN it passes judged since the parent
         2. THEN a plain commit of it fails
      4. WHEN a pull request splits code from its test plus doc across commits
         1. THEN it passes judged since the merge base
