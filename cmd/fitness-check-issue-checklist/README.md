# issue-checklist

Fails an issue body whose task list still has unchecked items, so closing it would close unfinished work.

It reads the body from a `--body-file` path, `-` meaning stdin, or from the runner's context-inline `--body` value. `fitness close-check` runs it on a closed issue, and `fitness pr-check` on each issue a pull request closes.

## Behavior

[Requirement 0001](docs/requirements/0001-closed-issues-have-no-open-tasks.md)

## Example

```
fitness issue-checklist
fitness issue-checklist --body "text"
```

It reads no `.fitnessrc.json` key.

