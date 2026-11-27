# semantic-commit

Requires the commit subject to follow [Conventional Commits](https://www.conventionalcommits.org/), `type(scope): description`, with the scope required. A merge commit (`Merge ...`) is always accepted.

## Allowed types

The types come from [conventional-commit-types](https://github.com/commitizen/conventional-commit-types), inlined in the check binary: the same eleven types, in the same order.

Scope is required, as in `feat(api): add endpoint`. Merge commits are always accepted.

## Behavior

[Requirement 0001](docs/requirements/0001-subjects-carry-type-and-scope.md)

## Where a full run starts

A full run finds the base the way the `local-branch` check does. A pull request run in GitHub Actions starts from the fetched `origin/$GITHUB_BASE_REF`. Locally, a branch counts the commits no other origin branch holds, whatever the default branch is called.

## The default branch

The default branch is `origin/HEAD`, else a fetched `origin/main` or `origin/master`. On it, on a detached HEAD, or while the default is unknown, only commits origin lacks count. With no commits of its own, the check passes with 0 checked.

A repo with no origin refs inherits nothing, so a full run judges the HEAD commit subject.

## Example

```
fitness semantic-commit
fitness semantic-commit --message "text"
```

It reads no `.fitnessrc.json` key.

