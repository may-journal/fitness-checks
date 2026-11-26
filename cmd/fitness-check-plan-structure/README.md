# plan-structure

Validates that a plan, a GitHub Issue body under the `Plan` label, follows the [Plan issue template](https://github.com/may-journal/fitness-runner/blob/main/.github/ISSUE_TEMPLATE/plan.md). It requires a one-line blockquote pitch, a `## Background` section, and a `## What needs to happen` section with at least one checklist item. No other `##` sections are allowed.

A repo's own `.github/ISSUE_TEMPLATE` file labeled `Plan` replaces the sections below. A section whose template shows a checklist item needs a checklist item.

Plans no longer live as repo files, so this check reads its target from stdin (or a `--body-file` path, `-` meaning stdin). It falls back to the runner's context-inline `--body` value. With no input at all it passes with zero files checked, so the file runner never trips on it. It is driven by the [issue-check workflow](https://github.com/may-journal/fitness-runner/blob/main/.github/workflows/README.md), not the local suite.

## Behavior

[Requirement 0001](docs/requirements/0001-plans-follow-the-template.md) states what this check does, and a test case named for each step proves it.
