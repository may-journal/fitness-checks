# comment-structure

Judges one GitHub issue or PR comment, so a reply stays short, readable, and free of agent noise. It holds a word budget outside tables, code, and diagrams, and allows at most one plain table. Diagrams are numbered, and change-log notes, agent frames, and model identifiers fail. It is driven by `fitness comment-check`, locally on a draft and in the [comment-check workflow](https://github.com/may-journal/fitness-runner/blob/main/.github/workflows/README.md#comment-check).

A comment is not a file in the tree. This check reads it from stdin, a `--body-file` path (`-` meaning stdin), or the runner's context-inline `--body` value. With no input it passes with zero files checked. Under `fitness --blob` it skips Plan, Epic, and PR bodies, which have their own structure checks.

## Behavior

[Requirement 0001](docs/requirements/0001-comments-stay-short-and-clean.md)

## Example

```
fitness comment-structure
fitness comment-structure --body "text"
```

Run `fitness explain comment-structure` for what each key does. This sets each key to its default:

```json
{
  "comment": {
    "maxWords": 100,
    "maxTables": 1
  }
}
```

