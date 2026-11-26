# 0001 Duplicate Code Stays Under One Percent

## Why

I get a repo whose copied lines stay within 1% of the lines it scans, judged by a built-in clone detector.

## Measurement

repos whose duplicated lines stay at or under 1%

-

repos with tracked text files

## Requirements

1. Shares past the limit fail
   1. GIVEN Developer has a repo with tracked text files that may share blocks of lines
      1. WHEN duplicated lines pass 1% of the scanned lines
         1. THEN it fails reporting `ERROR: duplicate code is X% of scanned lines, over the Y% threshold`
         2. THEN it reports `path: N duplicated lines` for each file holding them
      2. WHEN a pull request changes some files
         1. THEN it reads the changed files
         2. THEN it reads the tracked files sharing their endings
      3. WHEN Markdown or JSON files share a block
         1. THEN they count like any other tracked text file
      4. WHEN shared blocks stay under 5 lines or 50 tokens even with a table separator row in them
         1. THEN it passes
2. Excused content never counts
   1. GIVEN Developer has a repo whose tracked files may hold lines a copy check should not count
      1. WHEN lockfiles or build wrapper files or files marked `linguist-generated` share a block
         1. THEN they are left out
      2. WHEN a binary or untracked file shares a block
         1. THEN it is left out
      3. WHEN a block sits between the markers `duplicate-code:ignore-start` `duplicate-code:ignore-end`
         1. THEN it never counts
         2. THEN the deprecated marker `jscpd:ignore-start` works the same
         3. THEN the deprecated marker `jscpd:ignore-end` works the same
      4. WHEN docs repeat a doc template or a legend diagram or a style statement
         1. THEN the template is left out
         2. THEN those lines are left out
         3. THEN the other docs are still compared
