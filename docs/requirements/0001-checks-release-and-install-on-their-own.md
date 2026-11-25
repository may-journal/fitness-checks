# 0001 Checks Release And Install On Their Own

## Why

I take one check's fix without the others. Each check builds, tests, releases and installs alone, small and listed in one index.

## Measurement

Checks that build, test and release alone and appear in the release index

-

Checks

## Requirements

```gherkin
SCENARIO: 0001.1
  GIVEN Developer has every check folder in the repo
  WHEN each check is built on its own and its Go dependencies are listed
  THEN it holds no engine or other check or engine config
    AND it describes itself and its key defaults
    AND fitness runs it

SCENARIO: 0001.2
  GIVEN Developer changes a check or its tests or a shared helper or `go.mod`
  WHEN CI picks the checks to build and test for the change
  THEN it names only the checks the change can alter
    AND it builds and vets and tests each one

SCENARIO: 0001.3
  GIVEN Developer has merged commits of each kind to main and with and without earlier release tags
  WHEN the release workflow plans the releases
  THEN each changed check gets the bump its commits call for
    AND no other check releases

SCENARIO: 0001.4
  GIVEN Developer has a release tag for one check
  WHEN the release workflow builds that check
  THEN it writes a stripped and compressed archive and checksum for each platform named by tag

SCENARIO: 0001.5
  GIVEN Developer has releases at several majors for a check or with moved or missing major tags
  WHEN the release workflow moves the major tags
  THEN each major tag points at the newest release of its major
    AND a refused move fails with GitHub's message

SCENARIO: 0001.6
  GIVEN Developer has published and draft and prerelease releases of several checks
  WHEN the release workflow builds the release index
  THEN it lists each published check with the newest version of each major and its artifacts and checksums and describe fields
    AND it fetches only new versions

SCENARIO: 0001.7
  GIVEN Developer has a size baseline for the check programs
  WHEN the size check and the measure script run on the repo
  THEN a program over its limit fails by name
    AND the sizes are recorded

SCENARIO: 0001.8
  GIVEN Developer has each check's README
  WHEN the README test reads it
  THEN it holds the check's name and its summary and working links
```
