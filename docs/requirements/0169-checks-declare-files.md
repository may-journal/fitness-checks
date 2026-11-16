# 0169 Checks Declare Files

## Why

An engine runs a check only when a run reaches it. So each check states which files reach it and what else it requires. A pattern that does not compile, or a requirement no engine knows, would skip a check that should run.

## Measurement

Checks whose describe fields compile and reach the files they judge
-
Checks

## Requirements

- 0169.1
    - Given Developer has every check program this repo builds
        - When each prints `--describe`
            - Then its `files`, `exclude` and `needs` patterns compile, and it requires only what fitness-runner defines
- 0169.2
    - Given Developer has the gofmt, markdown-links, workflow-reporting and eslint checks
        - When their patterns meet a Go file, a vendored Go file, a Markdown doc, a workflow and a package.json
            - Then each reaches the files it judges and no other, and eslint needs the package.json
