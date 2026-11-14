# 0160 Checks engine free

## Why

Every check builds as its own program, without the engine or another check. So a check builds and releases on its own, and a repo installs only the checks it uses.

## Measurement

Checks whose Go dependencies hold no engine package and no other check
-
Checks

## Requirements

- 0160.1
    - Given Developer lists the Go dependencies of every check, tests included
        - When `go list -deps -test` runs on each check folder
            - Then none holds the engine's `cmd/fitness`, another check, or `metadata`, `aftereffect`, `dupe`, `report` or a package under them
- 0160.2
    - Given Developer builds every check folder
        - When each program runs with `--describe`
            - Then it prints what the pinned engine's bundled copy of that check prints, plus protocol 1
- 0160.3
    - Given Developer runs the requirement suite with `tests/requirements.sh`
        - When the installed engine runs every check
            - Then it starts this repo's program for each check, never its own built-in one
