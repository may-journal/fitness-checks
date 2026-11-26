# epic-structure

Validates that an epic follows the [Epic issue template](https://github.com/may-journal/fitness-runner/blob/main/.github/ISSUE_TEMPLATE/epic.md), or the repo's own template labeled `Epic`. An epic is a GitHub Issue body under the `Epic` label. Required: a one-line blockquote pitch, a `## Background` section, and a `## What needs to happen` list of issue references in run order. No other `##` sections are allowed.

Epics live as GitHub Issues, so this check reads its target from stdin (or a `--body-file` path, `-` meaning stdin). It falls back to the runner's context-inline `--body` value. With no input at all it passes with zero files checked. It is driven by `fitness issue-check` in the [issue-check workflow](https://github.com/may-journal/fitness-runner/blob/main/.github/workflows/README.md), not the local suite.

## Behavior

[Requirement 0001](docs/requirements/0001-epics-follow-the-epic-template.md) states what this check does, and a test case named for each step proves it.

