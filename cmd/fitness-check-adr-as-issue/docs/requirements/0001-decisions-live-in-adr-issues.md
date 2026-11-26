# 0001 Decisions Live In ADR Issues

## Why

I discuss, link and supersede each architecture decision in one issue, never in a file that drifts.

## Measurement

repos that keep no ADR in a file

-

repos checked

## Requirements

1. ADR files stay out of the repo
   1. GIVEN Developer has a repo
   2. GIVEN it holds Markdown files
   3. GIVEN it holds check READMEs
   4. GIVEN it may hold `adrAsIssue.allow` entries
      1. WHEN a Markdown file sits in a folder named `adr` in any case
         1. THEN it fails naming the file
         2. THEN it names the rule it broke
         3. THEN it names `fitness adr-migrate`
         4. THEN the folder's README passes
      2. WHEN a file's first top-level heading starts with the word ADR
         1. THEN it fails
         2. THEN a later ADR heading passes
         3. THEN a word such as `adr-checks` passes
      3. WHEN a file holds the three sections `## Context` `## Decision` `## Consequences`
         1. THEN it fails
         2. THEN the same sections in a code block pass
      4. WHEN a file only points at an issue
         1. THEN a file of one issue URL passes
         2. THEN a check README that links its issue passes
         3. THEN a file under `.github/` passes
         4. THEN it counts the paths it judged
      5. WHEN an allow entry names waiting ADR files with its reason
         1. THEN those files pass
         2. THEN every other ADR file still fails
      6. WHEN an allow entry is invalid
         1. THEN an entry that lacks a reason fails naming its path
         2. THEN an entry that is no glob fails naming its path
         3. THEN an entry that matches no ADR file fails naming its path
      7. WHEN git cannot list the repo's files
         1. THEN it fails
