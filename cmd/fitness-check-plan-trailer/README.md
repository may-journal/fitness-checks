# plan-trailer

Validates an optional `Plan #NN` commit trailer. A commit may name the Plan Issue it implements with a line reading exactly `Plan #<number>`. The trailer is optional, so a message without one passes. A plan reference that is present must be well formed.

A Plan in another repo of the same owner, such as an org-wide Plan in the org's `.github` repo, reads `Plan <owner>/<repo>#<number>` or `Plan [https://github.com/](https://github.com/)<owner>/<repo>/issues/<number>`. The trailer goes on the branch's first commit only, since [issue-link-once](../fitness-check-issue-link-once/README.md) fails a repeat in any of the spellings.

Like `semantic-commit`, it reads the proposed message from the context-inline `--message` value, falling back to the HEAD commit message. The `commit-msg` hook runs it beside `semantic-commit`.

## Behavior

[Requirement 0001](docs/requirements/0001-plan-trailers-are-spelled-right.md) states what this check does, and a test case named for each step proves it.

The pre-push gate accepts a qualified Plan only from the pushing repo's own owner, and reads it from the repo it names.
