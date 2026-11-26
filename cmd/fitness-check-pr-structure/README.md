# pr-structure

Validates that a pull request description follows the [PR template](https://github.com/may-journal/fitness-runner/blob/main/.github/PULL_REQUEST_TEMPLATE.md). It requires a one-line blockquote summary, a `## Background` section, and a `## Changelog` section with a bullet, and no other `##` sections. The "what changed" lives in the diff, so the description carries only the why and the changelog.

A PR description is not a file in the tree, so this check reads its target from stdin. It also accepts a `--body-file` path (`-` meaning stdin) and falls back to the runner's context-inline `--body` value. With no input at all it passes with zero files checked, so the file runner never trips on it. It is driven by the [pr-check workflow](https://github.com/may-journal/fitness-runner/blob/main/.github/workflows/README.md), not the local suite.

## Behavior

[Requirement 0001](docs/requirements/0001-pr-bodies-follow-the-template.md)

## Notes

The bullet links the PR's [entry file](../fitness-check-changelog/README.md#entry-files) rather than copying it. This check sees only the body, so `pr-check` enforces the link with the PR's diff. It fails a PR that adds an entry without linking it, links another file, or copies a bullet the PR adds.

## Example

```
fitness pr-structure
fitness pr-structure --body "text"
```

It reads no `.fitnessrc.json` key.

