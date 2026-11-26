# 0007 The Repo Holds Only Working Source

## Why

I review real source only. Plan files, copied code, repeated values, stale pins, unexplained ignore lines and broken workflows never land in the repo.

## Measurement

repo files and settings that pass every working-source rule

-

repo files and settings checked

## Requirements

1. Removed in 2026.11.25.0800: moved to the `build-output-untracked` doc
2. Design scraps stay out
   1. GIVEN Developer has plan files and mockups
      1. WHEN the source checks look for temporary design work
         1. THEN each such file fails naming its rule
         2. THEN maintained docs pass
3. Removed in 2026.11.26.0100: moved to the duplicate-code doc
4. Repeated strings get names
   1. GIVEN Developer has string values used several times
      1. WHEN the source checks count repeated values
         1. THEN a value used three or more times fails with its count and places
5. Removed in 2026.11.26.0100: moved to the dependency-currency doc
6. Removed in 2026.11.26.0100: moved to the gitignore-why doc
7. Workflow calls resolve
   1. GIVEN Developer has workflow `uses:` lines at tags and shas
      1. WHEN the source checks look up each pinned path
         1. THEN a path missing at its ref fails naming file and line and ref
         2. THEN other calls pass
8. Workflow jobs stay clean
   1. GIVEN Developer has workflow jobs with or without a summary step
      1. WHEN the source checks read each workflow job
         1. THEN a job with no `always()` summary fails naming the line
         2. THEN an inline script or invalid YAML fails naming the line
