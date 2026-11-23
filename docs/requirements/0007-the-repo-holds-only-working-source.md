# 0007 The Repo Holds Only Working Source

## Why

I review real source only. Build output, plan files, copied code, repeated values, stale pins, unexplained ignore lines and broken workflows never land in the repo.

## Measurement

repo files and settings that pass every working-source rule

-

repo files and settings checked

## Requirements

- 0007.1
  - Given Developer has tracked, ignored and imported `dist`, Maven and Gradle output
    - When the source checks look for build output
      - Then tracked, unignored or imported output fails with a fix
- 0007.2
  - Given Developer has plan files, mockups and review artifacts in folders, names and checklists
    - When the source checks look for temporary design work
      - Then each such file fails naming its rule, and maintained docs pass
- 0007.3
  - Given Developer has tracked text files that share blocks of lines
    - When the source checks measure duplicated code
      - Then a share past the limit fails naming each file, and generated or marked blocks pass
- 0007.4
  - Given Developer has string values used twice, three times or more across source and test files
    - When the source checks count repeated values
      - Then a value used three or more times fails with its count and places
- 0007.5
  - Given Developer has declared dependencies at, behind or missing from their latest versions
    - When the source checks compare the pins with the registry
      - Then each stale or missing one fails as `name: current → latest`, and an unreachable registry only adds info
- 0007.6
  - Given Developer has `.gitignore` patterns with and without a comment right above
    - When the source checks read the ignore file
      - Then each pattern without its own non-empty comment fails naming the line, and a repo with no file passes
- 0007.7
  - Given Developer has workflow `uses:` lines for may-journal paths at tags and shas
    - When the source checks look up each pinned path
      - Then a path missing at its ref fails naming file, line and ref, and other calls pass
- 0007.8
  - Given Developer has workflow jobs with or without a summary step, and inline shell in steps
    - When the source checks read each workflow job
      - Then a job with no `always()` summary step, an inline script, or invalid YAML fails naming its line
