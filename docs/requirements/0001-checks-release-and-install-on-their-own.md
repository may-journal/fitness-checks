# 0001 Checks Release And Install On Their Own

## Why

I take one check's fix without the others. Each check builds, tests, releases and installs alone, small and listed in one index.

## Measurement

Checks that build, test and release alone and appear in the release index

-

Checks

## Requirements

1. Checks stand alone
   1. GIVEN Developer has every check folder
      1. WHEN each check is built alone
         1. THEN it holds no engine or other check's code
         2. THEN it describes itself and declares key defaults
         3. THEN fitness runs it
2. CI tests changes
   1. GIVEN Developer changes a check
   2. GIVEN its tests change
   3. GIVEN a shared helper changes
   4. GIVEN `go.mod` changes
      1. WHEN CI picks the checks to test
         1. THEN it tests only checks the change can alter
      2. WHEN CI splits the suite into jobs
         1. THEN each root test runs once
         2. THEN only changed checks' folders run once in six shares
3. Releases follow commits
   1. GIVEN Developer has commits with and without release tags
      1. WHEN the release workflow plans the releases
         1. THEN each changed check gets the bump its commits call for
         2. THEN no other check releases
4. Archives ship per platform
   1. GIVEN Developer has a release tag for one check
      1. WHEN the release workflow builds that check
         1. THEN it writes a stripped archive and checksum per platform
5. Major tags follow releases
   1. GIVEN Developer has releases at several majors with moved or missing major tags
      1. WHEN the release workflow moves the major tags
         1. THEN each points at its newest release
         2. THEN a refused move fails
6. The index lists new versions
   1. GIVEN Developer has published and draft releases
      1. WHEN the release workflow builds the release index
         1. THEN it lists each published check's newest version per major
         2. THEN each entry holds artifacts and checksums and describe fields
         3. THEN it fetches only new versions
7. Sizes stay limited
   1. GIVEN Developer has a size baseline for the programs
      1. WHEN the size check runs
         1. THEN a program over its limit fails by name
         2. THEN the sizes are recorded
8. READMEs name checks
   1. GIVEN Developer has each check's README
      1. WHEN the README test reads it
         1. THEN it holds its name and summary and links
