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
- 0169.3
    - Given Developer has the PR, issue, comment and commit message checks
        - When their patterns meet each virtual body path, piped text, the commit message path and a Go file
            - Then each body check reaches its own body and piped text, each message check its message, and none a Go file
- 0169.4
    - Given Developer has the prose checks
        - When their patterns meet a Markdown doc, a PR body, a comment, piped text, a stylesheet and a Go file
            - Then the body readers reach the Markdown and bodies, the comment readers reach the Go file, and none reaches the stylesheet
