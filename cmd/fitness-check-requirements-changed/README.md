# requirements-changed

Fails a change that adds or changes product code without also changing a test and a requirement doc. The [requirements](../fitness-check-requirements/README.md) check judges the docs and tests that exist; this check judges what a commit or pull request leaves out.

## What counts

- Product code: `.go` and `.swift` files that are not tests, and `.sh` scripts.
- Vendored code is not product code: files under `vendor`, `third_party`, `Pods`, `Carthage`, or `node_modules`.
- Fixtures and build output are not either: Go `testdata`, `dist`, `.build`, `DerivedData`, and what each [language module](https://github.com/may-journal/fitness-runner/blob/main/go/pkg/lang/langs/langs.go) names.
- Generated files are not: those `.gitattributes` marks `linguist-generated`, and Go files that say `Code generated … DO NOT EDIT.`
- Synced files: a file `.github/synced-files.json` names passes when its staged blob matches, and fails naming its source repo otherwise.
- Tests: `_test.go` files, and Swift files that declare a test or sit in a `…Tests` folder, as the requirements check pairs them.
- Requirement docs: `docs/requirements/NNNN-*.md`.

A file counts as changed only when its lines change beyond whitespace, so a pure rename or a reformat changes nothing. Nor does a file whose code, with comments removed, is the same before and after. That holds for tests and docs too. Directives, a shebang, a cgo preamble, and a Go example's output comment count as code.

## Behavior

[Requirements](docs/requirements)

## When it runs

- Commits: the staged files. With a `commit-msg` hook, the hook judges them.
- A `commit-msg` hook, as `fitness init` installs, also makes `pre-commit` leave the check out.
- Amends: everything the amended commit changes since its parent.
- Pull requests: CI judges the diff against the merge base with the base branch.

## Example

```
fitness requirements-changed
fitness requirements-changed --message "text"
```

It reads no `.fitnessrc.json` key.

