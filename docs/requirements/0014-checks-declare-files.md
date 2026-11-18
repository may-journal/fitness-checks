# 0014 Checks Declare Files

## Why

An engine runs a check only when a run reaches it. So each check states which files reach it and what else it requires. A pattern that does not compile, or a requirement no engine knows, would skip a check that should run.

## Measurement

Checks whose describe fields compile and reach the files they judge
-
Checks

## Requirements

- 0014.1
    - Given Developer has every check program this repo builds
        - When each prints `--describe`
            - Then its `files`, `exclude` and `needs` patterns compile, and it requires only what fitness-runner defines
- 0014.2
    - Given Developer has the gofmt, markdown-links, workflow-reporting and eslint checks
        - When their patterns meet a Go file, a vendored Go file, a Markdown doc, a workflow and a package.json
            - Then each reaches the files it judges and no other, and eslint needs the package.json
- 0014.3
    - Given Developer has the PR, issue, comment and commit message checks
        - When their patterns meet each virtual body path, piped text, the commit message path and a Go file
            - Then a body check reaches its body and piped text, a message check its message, neither a Go file, and each needs fitness 1.75.0
- 0014.4
    - Given Developer has the prose checks
        - When their patterns meet a Markdown doc, each body, piped text, a stylesheet, a Go file and a script with no extension
            - Then the body readers reach the Markdown and bodies, the comment readers reach the Go file and script, and none reaches the stylesheet
- 0014.5
    - Given Developer has a describe line whose `files` pattern does not compile, and a check's own describe line
        - When `scripts/describe-patterns`, which `scripts/release-build.sh` runs before a release, reads each
            - Then it exits 1 naming the release, field and pattern of the first, and exits 0 for the second
