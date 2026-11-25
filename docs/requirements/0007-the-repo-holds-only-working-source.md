# 0007 The Repo Holds Only Working Source

## Why

I review real source only. Build output, plan files, copied code, repeated values, stale pins, unexplained ignore lines and broken workflows never land in the repo.

## Measurement

repo files and settings that pass every working-source rule

-

repo files and settings checked

## Requirements

```gherkin
SCENARIO: 0007.1
  GIVEN Developer has tracked and ignored and imported `dist` and Maven and Gradle output
  WHEN the source checks look for build output
  THEN tracked or unignored or imported output fails with a fix

SCENARIO: 0007.2
  GIVEN Developer has plan files and mockups and review artifacts in folders and names and checklists
  WHEN the source checks look for temporary design work
  THEN each such file fails naming its rule
    AND maintained docs pass

SCENARIO: 0007.3
  GIVEN Developer has tracked text files that share blocks of lines
  WHEN the source checks measure duplicated code
  THEN a share past the limit fails naming each file
    AND generated or marked blocks pass

SCENARIO: 0007.4
  GIVEN Developer has string values used twice or three times or more across source and test files
  WHEN the source checks count repeated values
  THEN a value used three or more times fails with its count and places

SCENARIO: 0007.5
  GIVEN Developer has declared dependencies at or behind or missing from their latest versions
  WHEN the source checks compare the pins with the registry
  THEN each stale or missing one fails as `name: current → latest`
    AND an unreachable registry only adds info

SCENARIO: 0007.6
  GIVEN Developer has `.gitignore` patterns with and without a comment right above
  WHEN the source checks read the ignore file
  THEN each pattern without its own non-empty comment fails naming the line
    AND a repo with no file passes

SCENARIO: 0007.7
  GIVEN Developer has workflow `uses:` lines for may-journal paths at tags and shas
  WHEN the source checks look up each pinned path
  THEN a path missing at its ref fails naming file and line and ref
    AND other calls pass

SCENARIO: 0007.8
  GIVEN Developer has workflow jobs with or without a summary step and inline shell in steps
  WHEN the source checks read each workflow job
  THEN a job with no `always()` summary step or an inline script or invalid YAML fails naming the line
```
