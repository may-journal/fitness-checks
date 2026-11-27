# 0001 Product Code Changes With Its Test And Doc

## Why

I get a change to product code held until a test and a requirement doc change with it.

## Measurement

changes to product code that also change a test and a requirement doc

-

changes to product code

## Requirements

1. Product code is paired
   1. GIVEN Developer has a repo with Go or Swift code
   2. GIVEN it has requirement docs
      1. WHEN changed product code lacks a changed test or a changed requirement doc
         1. THEN it fails naming the file
         2. THEN it names the missing half
      2. WHEN product code changes beside a changed test plus a changed requirement doc
         1. THEN it passes
      3. WHEN only docs or config or the changelog change or a file is renamed or reformatted or only its comments change
         1. THEN it passes since nothing there needs a pair
      4. WHEN a merge or a cherry-pick or a revert replays product code
         1. THEN it is not judged again
         2. THEN its pull request diff still is
      5. WHEN product code is deleted
         1. THEN it passes if no doc names it
         2. THEN it fails naming the file with the doc if a doc names it
         3. THEN it names the tests of that doc
      6. WHEN the repo has no requirement doc yet
         1. THEN it passes since there is no doc to pair with
      7. WHEN a commit message asks to skip the check
         1. THEN it still fails since no trailer or flag or setting turns it off
2. Other files need no pair
   1. GIVEN Developer has a repo with Go or Swift code
   2. GIVEN it has requirement docs
      1. WHEN only vendored files or `testdata` or build output such as `dist` or `.build` or `DerivedData` or generated files change
         1. THEN it passes since they are not product code
      2. WHEN a full scan runs with nothing staged
         1. THEN it judges nothing
         2. THEN it passes
