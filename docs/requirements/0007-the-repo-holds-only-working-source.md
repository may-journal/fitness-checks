# 0007 The Repo Holds Only Working Source

## Why

I review real source only. Plan files, copied code, repeated values, stale pins, unexplained ignore lines and broken workflows never land in the repo.

## Measurement

repo files and settings that pass every working-source rule

-

repo files and settings checked

## Requirements

1. Removed in 2026.11.25.0800: moved to the `build-output-untracked` doc
2. Removed in 2026.11.26.0300: moved to the no-plans-dir doc
3. Removed in 2026.11.26.0100: moved to the duplicate-code doc
4. Repeated strings get names
   1. GIVEN Developer has string values used several times
      1. WHEN the source checks count repeated values
         1. THEN a value used three or more times fails with its count and places
5. Removed in 2026.11.26.0100: moved to the dependency-currency doc
6. Removed in 2026.11.26.0100: moved to the gitignore-why doc
7. Removed in 2026.11.26.0200: moved to the workflow-ref-exists doc
8. Removed in 2026.11.26.0200: moved to the workflow-reporting and workflow-scripts docs
