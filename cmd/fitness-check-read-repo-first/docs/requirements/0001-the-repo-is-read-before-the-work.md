# 0001 The Repo Is Read Before The Work

## Why

I see a reminder to read the repo's decisions and enabled checks before I work in it.

## Measurement

runs that show the reminder with their enabled checks

-

runs of `read-repo-first`

## Requirements

1. The reminder shows with the enabled checks
   1. GIVEN Developer has a run that enables some checks
      1. WHEN fitness runs `read-repo-first`
         1. THEN it passes with zero files
         2. THEN it asks whether the decisions were read
         3. THEN it asks whether the enabled checks were read
         4. THEN it warns never to use `--no-verify`
      2. WHEN a run enables `duplicate-code` without `gofmt`
         1. THEN the table lists only the enabled checks
      3. WHEN a check name or path is too long for its cell
         1. THEN the cell is shortened with an ellipsis
      4. WHEN color stays enabled
         1. THEN the table is colored
      5. WHEN color is disabled with `NO_COLOR`
         1. THEN the table prints plain
