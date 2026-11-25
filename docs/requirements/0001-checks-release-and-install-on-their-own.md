# 0001 Checks Release And Install On Their Own

## Why

I take one check's fix without the others. Each check builds, tests, releases and installs alone, small and listed in one index.

## Measurement

Checks that build, test and release alone and appear in the release index

-

Checks

## Requirements

- 0001.1
  - Given Developer has every check folder in the repo
    - When each check is built on its own and its Go dependencies are listed
      - Then it holds no engine, other check or engine config, describes itself, key defaults, and fitness runs it
- 0001.2
  - Given Developer changes a check or its tests or a shared helper or `go.mod`
    1. When CI picks the checks to build and test for the change
       - Then it names only the checks the change can alter and builds and tests each one
    2. When CI splits the requirement suite into parallel jobs for the change
       - Then each root test runs in exactly one job and only the changed checks' folders run
- 0001.3
  - Given Developer has merged commits of each kind to main, with and without earlier release tags
    - When the release workflow plans the releases
      - Then each changed check gets the bump its commits call for, and no other check releases
- 0001.4
  - Given Developer has a release tag for one check
    - When the release workflow builds that check
      - Then it writes a stripped, compressed archive and checksum for each platform, named by tag
- 0001.5
  - Given Developer has releases at several majors for a check, with moved or missing major tags
    - When the release workflow moves the major tags
      - Then each major tag points at the newest release of its major, and a refused move fails with GitHub's message
- 0001.6
  - Given Developer has published, draft and prerelease releases of several checks
    - When the release workflow builds the release index
      - Then it lists each published check with the newest version of each major, its artifacts, checksums and describe fields, fetching only new versions
- 0001.7
  - Given Developer has a size baseline for the check programs
    - When the size check and the measure script run on the repo
      - Then a program over its limit fails by name, and the sizes are recorded
- 0001.8
  - Given Developer has each check's README
    - When the README test reads it
      - Then it holds the check's name, its summary and working links
