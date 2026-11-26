# changelog-updated

Checks that each staged commit writes a changelog entry about its change. It runs when the runner provides staged file context. The added lines in the branch's [entry file](https://github.com/may-journal/fitness-checks/blob/main/cmd/fitness-check-changelog/README.md#entry-files) must share at least three words with the rest of the staged diff.

## Behavior

[Requirement 0001](docs/requirements/0001-commits-describe-their-change-in-an-entry-file.md)

## Example

```
fitness changelog-updated
```

It reads no `.fitnessrc.json` key.

