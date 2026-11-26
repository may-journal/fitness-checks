# 0001 Changelog Bullets Are Few, Short And Typed

## Why

I read each change's entry at a glance: a few short bullets, each naming its type.

## Measurement

entry files whose sections pass every bullet rule

-

entry files

## Requirements

```gherkin
SCENARIO: 0001.1
  GIVEN Developer has a repo with or without `changelog/` entry files
  WHEN each section holds 3 to 5 bullets under 365 characters
    AND each bullet opens with a semantic type
  THEN it passes
    AND it counts the entry files
    AND types may read `Feat:` or `feat(scope)!:`
  WHEN a `###` section holds fewer than 3 `- ` bullets
    AND a `###` section holds more than 5 `- ` bullets
  THEN it fails naming the file
    AND it names the section
    AND it names its count
  WHEN a bullet reaches 365 characters with its continuation lines joined
  THEN it fails naming the file
    AND it names the line
    AND it names the length
  WHEN a bullet opens with none of the types from `Feat:` to `Revert:`
    AND the bullet is no such commit subject
  THEN it fails that bullet alone naming its line
  WHEN several entry files break these rules
  THEN each error names its own file
  WHEN the repo has no entry files
    AND the repo has only a `CHANGELOG.md`
  THEN it passes with no files
```
