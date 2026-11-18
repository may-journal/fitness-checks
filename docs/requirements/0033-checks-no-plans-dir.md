# 0033 Checks No Plans Dir

## Why

I keep plans in `Plan` issues and mockups or reviews in comments. Temporary work never lands in the repo unnoticed.

## Measurement

plan files and temporary design artifacts reported red
-
plan files and temporary design artifacts in the repo

## Requirements

- 0033.1
    - Given Developer has a plan file in `docs/plans/`
        - When fitness runs no-plans-dir
            - Then it fails naming the file and telling me to open a Plan Issue
- 0033.2
    - Given Developer has a plan file in a folder below `docs/plans/`
        - When fitness runs no-plans-dir
            - Then it fails naming that nested file too
- 0033.3
    - Given Developer has maintained docs, a shipped image named `screen.png`, and a test that quotes plan text inline
        - When fitness runs no-plans-dir
            - Then it passes and leaves every file in place
- 0033.4
    - Given Developer has files in a `Plans`, `mockups`, `validation` or `review-artifacts` folder anywhere in the repo
        - When fitness runs no-plans-dir
            - Then it fails naming each file, its folder rule and the comment to move it to, sorted
- 0033.5
    - Given Developer has files named `*-plan.md`, `*-mockup.*` and `pr-<number>-*`
        - When fitness runs no-plans-dir
            - Then it fails naming each file and its rule
- 0033.6
    - Given Developer has a Markdown file holding a `## What needs to happen` checklist outside `.github/`
        - When fitness runs no-plans-dir
            - Then it fails naming the file, while the same checklist in a `.github/` template passes
- 0033.7
    - Given Developer has a plan file staged but not committed, or renamed into `mockups/`
        - When fitness runs no-plans-dir
            - Then it fails naming the staged path or the rename's target
- 0033.8
    - Given Developer has a repo whose git index cannot be read
        - When fitness runs no-plans-dir
            - Then it fails saying it cannot list the repo's files
