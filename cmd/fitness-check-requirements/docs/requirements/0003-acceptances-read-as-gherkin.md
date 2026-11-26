# 0003 Acceptances Read As Gherkin

## Why

I number each acceptance and give it a name. Under it I write numbered GIVEN, WHEN and THEN items in caps, one condition or result each. The standard Gherkin parser reads them.

## Measurement

Numbered acceptances the check reads and pairs with their tests

-

Numbered acceptances

## Requirements

1. Caps keywords
   1. GIVEN Developer has numbered acceptances `requirements` reads
      1. WHEN each level numbers from 1
         1. THEN it pairs steps with tests
      2. WHEN a keyword is not in caps
         1. THEN it names the keyword
      3. WHEN an item is an AND or BUT
         1. THEN it asks for the next GIVEN or THEN
      4. WHEN an item holds a comma outside code
         1. THEN it asks to split
      5. WHEN an item is misnumbered or misnamed
         1. THEN it names the need
      6. WHEN the parser rejects an item
         1. THEN it fails at that line
2. Removed in 2026.11.25.1830: the check no longer converts bullet docs
3. Tests need no change
   1. GIVEN Developer moves a bullet doc to numbered items
      1. WHEN `requirements` judges it
         1. THEN its tests need no change
4. Renames and moves
   1. GIVEN Developer has numbered acceptances
      1. WHEN one is renamed in place
         1. THEN it passes
      2. WHEN a commit moves one
         1. THEN it names both IDs
      3. WHEN it moves back
         1. THEN it passes
5. Removed lines
   1. GIVEN Developer has numbered acceptances
      1. WHEN a commit removes one with its test
         1. THEN it passes
      2. WHEN a `Removed` line holds items or lacks a version or reason
         1. THEN it fails
      3. WHEN a test names a `Removed` acceptance
         1. THEN it fails
      4. WHEN it is restored
         1. THEN it fails as reused
      5. WHEN one commit removes one then appends it
         1. THEN it passes
6. Removed lines uncapped
   1. GIVEN Developer has a capped doc with `Removed` lines
      1. WHEN either `requirements` or `prose-budget` reads it
         1. THEN it passes
7. Unknown flags fail
   1. GIVEN Developer runs `requirements` with a flag it does not have
      1. WHEN the flag is a removed `--up` or `--down`
         1. THEN it fails naming the flag
8. Last acceptance kept
   1. GIVEN Developer has numbered acceptances
      1. WHEN the last becomes `Removed`
         1. THEN it passes
      2. WHEN the last is deleted outright
         1. THEN it names the `Removed` line
9. Removed in 2026.11.26.2000: moved to 0008
