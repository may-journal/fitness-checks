# 0001 Duplicate Code Stays Under One Percent

## Why

I get a repo whose copied lines stay within 1% of the lines it scans, judged by a clone detector built into the check.

## Measurement

repos whose duplicated lines stay at or under 1%

-

repos with tracked text files

## Requirements

```gherkin
SCENARIO: 0001.1
  GIVEN Developer has a repo with tracked text files that may share blocks of lines
  WHEN duplicated lines pass 1% of the scanned lines
  THEN it fails reporting `ERROR: duplicate code is X% of scanned lines, over the Y% threshold` and `path: N duplicated lines` for each file holding them
  WHEN shared blocks stay under 5 lines or 50 tokens
  THEN it passes
  WHEN Markdown and JSON files share a block
  THEN they count like any other tracked text file
  WHEN a pull request changes some files
  THEN it reads the changed files and the tracked files sharing their endings
  WHEN a file holds commented-out code
  THEN that code never counts as a copy

SCENARIO: 0001.2
  GIVEN Developer has a repo whose tracked files may hold lines a copy check should not count
  WHEN a binary or untracked file shares a block
  THEN it is left out
  WHEN lockfiles or build wrapper files or files marked `linguist-generated` share a block
  THEN they are left out
  WHEN a block sits between `duplicate-code:ignore-start` and `duplicate-code:ignore-end` markers
  THEN it never counts and the deprecated `jscpd:ignore-start` and `jscpd:ignore-end` markers work the same
  WHEN Markdown repeats a legend diagram or a style statement in a mermaid block or a table separator row
  THEN those lines never count
  WHEN docs repeat a doc template
  THEN the template is left out and the other docs are still compared
```
