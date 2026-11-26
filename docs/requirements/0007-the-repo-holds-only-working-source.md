# 0007 The Repo Holds Only Working Source

## Why

I review real source only. Build output, plan files, copied code, repeated values, stale pins, unexplained ignore lines and broken workflows never land in the repo.

## Measurement

repo files and settings that pass every working-source rule

-

repo files and settings checked

## Requirements

1. Build output stays out
   1. GIVEN Developer has tracked `dist` or Maven or Gradle output
   2. GIVEN it has unignored or imported output
      1. WHEN the source checks look for build output
         1. THEN each fails with a fix
2. Design scraps stay out
   1. GIVEN Developer has plan files and mockups
      1. WHEN the source checks look for temporary design work
         1. THEN each such file fails naming its rule
         2. THEN maintained docs pass
3. Duplicates stay small
   1. GIVEN Developer has text files that share blocks of lines
      1. WHEN the source checks measure duplicated code
         1. THEN a share past the limit fails naming each file
         2. THEN generated or marked blocks pass
4. Repeated strings get names
   1. GIVEN Developer has string values used several times
      1. WHEN the source checks count repeated values
         1. THEN a value used three or more times fails with its count and places
5. Dependencies stay current
   1. GIVEN Developer has dependencies behind their latest versions
      1. WHEN the source checks compare the pins with the registry
         1. THEN a stale or missing one fails as `name: current → latest`
         2. THEN an unreachable registry only adds info
6. Ignores say why
   1. GIVEN Developer has `.gitignore` patterns with and without a comment right above
      1. WHEN the source checks read the ignore file
         1. THEN a pattern without its own comment fails naming the line
         2. THEN a repo with no file passes
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
