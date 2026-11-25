# 0001 Changelog Bullets Are Few, Short And Typed

## Why

I read each change's entry at a glance: a few short bullets, each naming its type.

## Measurement

entry files whose sections pass every bullet rule

-

entry files

## Requirements

- 0001.1
  - Given Developer has a repo, with or without `changelog/` entry files
    1. When each section holds 3 to 5 bullets under 365 characters, each opening with a semantic type
       - Then it passes, counting the entry files, whether types read `Feat:` or `feat(scope)!:`
    2. When a `###` section holds fewer than 3 or more than 5 `- ` bullets
       - Then it fails naming the file, the section and its count
    3. When a bullet, its continuation lines joined, reaches 365 characters
       - Then it fails naming the file, the line and the length
    4. When a bullet opens with none of the types `Feat:`, `Fix:`, ... `Revert:`, nor such a commit subject
       - Then it fails that bullet alone, naming its line
    5. When several entry files break these rules
       - Then each error names its own file
    6. When the repo has no entry files, only a `CHANGELOG.md`
       - Then it passes with no files
