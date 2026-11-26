# 0001 Decisions Live In ADR Issues

## Why

I discuss, link and supersede each architecture decision in one issue, never in a file that drifts.

## Measurement

repos that keep no ADR in a file

-

repos checked

## Requirements

```gherkin
SCENARIO: 0001.1
  GIVEN Developer has Markdown files AND check READMEs AND `adrAsIssue.allow` entries in any repo
  WHEN a Markdown file sits in a folder named `adr` in any case
  THEN it fails naming the file
    AND it names its rule
    AND it names `fitness adr-migrate`
    BUT the folder's README passes
  WHEN a file's first top-level heading starts with the word ADR
  THEN it fails
    BUT a later ADR heading passes
    AND a word such as `adr-checks` passes
  WHEN a file holds `## Context` AND `## Decision` AND `## Consequences` sections
  THEN it fails
    BUT the same sections in a code block pass
  WHEN a file is one issue URL
    AND a check README links its issue
    AND a file sits under `.github/`
  THEN it passes
    AND it counts the paths it judged
  WHEN an allow entry names waiting ADR files with its reason
  THEN those files pass
    BUT every other ADR file still fails
  WHEN an allow entry lacks a reason
    AND an entry is no glob
    AND an entry matches no ADR file
  THEN each such entry fails naming its path
  WHEN git cannot list the repo's files
  THEN it fails
```
