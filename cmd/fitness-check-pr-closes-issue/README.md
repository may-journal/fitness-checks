# pr-closes-issue

Validates that a pull request description closes at least one issue on merge ([#69](https://github.com/may-journal/fitness-runner/issues/69)). GitHub auto-closes a linked issue only for a fixed set of keywords. A PR that merely references an issue leaves it open, and this check fails that.

It runs in [pr-check.yml](https://github.com/may-journal/fitness-runner/blob/main/.github/workflows/pr-check.yml) beside `pr-structure`. It reads the PR body from stdin, a `--body-file` path (`-` meaning stdin), or the context-inline `--body`, and passes inert with no input. For a closed Plan, `fitness pr-check` reads that Plan's own closing targets via `--emit-closed`, and requires the PR to close them too. That flag prints the body's closed issue numbers as a JSON array.

## Behavior

[Requirement 0001](docs/requirements/0001-a-pr-closes-its-issue.md) states what this check does, and a test case named for each step proves it.

## Unchecked items

A PR must not close an issue whose checklist is unfinished ([#141](https://github.com/may-journal/fitness-runner/issues/141)). For each issue the PR closes, `fitness pr-check` lists every unchecked `- [ ]` item, skipping fenced code. Each one fails the PR, named with its issue number.

Tick each item as it is delivered, or stop closing the issue. After merge, `fitness close-check` reopens any issue that still closed with unchecked items.

## Closing keywords

The closing-keyword list is GitHub's fixed set; keep it in sync in `gitx.ClosingKeywords` if GitHub's set ever changes.
