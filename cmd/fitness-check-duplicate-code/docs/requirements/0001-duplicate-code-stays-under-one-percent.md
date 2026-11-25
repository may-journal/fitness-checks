# 0001 Duplicate Code Stays Under One Percent

## Why

I get a repo whose copied lines stay within 1% of the lines it scans, judged by a clone detector built into the check.

## Measurement

repos whose duplicated lines stay at or under 1%

-

repos with tracked text files

## Requirements

- 0001.1
  - Given Developer has a repo with tracked text files that may share blocks of lines
    1. When duplicated lines pass 1% of the scanned lines
       - Then it fails reporting `ERROR: duplicate code is X% of scanned lines, over the Y% threshold` and `path: N duplicated lines` for each file holding them
    2. When shared blocks stay under 5 lines or 50 tokens
       - Then it passes
    3. When Markdown and JSON files share a block
       - Then they count like any other tracked text file
    4. When a pull request changes some files
       - Then it reads the changed files and the tracked files sharing their endings
    5. When a file holds commented-out code
       - Then that code never counts as a copy
- 0001.2
  - Given Developer has a repo whose tracked files may hold lines a copy check should not count
    1. When a binary or untracked file shares a block
       - Then it is left out
    2. When lockfiles or build wrapper files or files marked `linguist-generated` share a block
       - Then they are left out
    3. When a block sits between `duplicate-code:ignore-start` and `duplicate-code:ignore-end` markers
       - Then it never counts and the deprecated `jscpd:ignore-start` and `jscpd:ignore-end` markers work the same
    4. When Markdown repeats a legend diagram or a style statement in a mermaid block or a table separator row
       - Then those lines never count
    5. When docs repeat a doc template
       - Then the template is left out and the other docs are still compared
