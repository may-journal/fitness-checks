# 0001 Decisions Live In ADR Issues

## Why

I discuss, link and supersede each architecture decision in one issue, never in a file that drifts.

## Measurement

repos that keep no ADR in a file

-

repos checked

## Requirements

- 0001.1
  - Given Developer has Markdown files, check READMEs and `adrAsIssue.allow` entries in any repo
    1. When a Markdown file sits in a folder named `adr`, in any case
       - Then it fails naming the file, its rule and `fitness adr-migrate`, and the folder's README passes
    2. When a file's first top-level heading starts with the word ADR
       - Then it fails, and a later ADR heading or a word such as `adr-checks` passes
    3. When a file holds `## Context`, `## Decision` and `## Consequences` sections
       - Then it fails, and the same sections in a code block pass
    4. When a file is one issue URL, a check README links its issue, or a file sits under `.github/`
       - Then it passes, counting the paths it judged
    5. When an allow entry names waiting ADR files with its reason
       - Then those files pass, and every other ADR file still fails
    6. When an allow entry lacks a reason, is no glob, or matches no ADR file
       - Then each such entry fails naming its path
    7. When git cannot list the repo's files
       - Then it fails
