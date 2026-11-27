# requirements

Enforces docs/requirements in every repo a language module covers. In a repo written in any language of the [module registry](https://github.com/may-journal/fitness-runner/blob/main/go/pkg/lang/langs/langs.go), each file in `docs/requirements` is one requirement.

A requirement says why the requirement matters, how to measure it, and which acceptances prove it. Each acceptance is owned by exactly one test named for it, and each test owns exactly one acceptance. A test that proves no requirement is either a requirement nobody wrote down or a duplicate to delete.

## Template

Copy this into `docs/requirements/NNNN-kebab-title.md`. Each rule is under Rules.

Write it as Biome formats Markdown: blank lines around the `-`, so it never underlines the numerator as a heading. Nest by three spaces.

```
# 0001 Offline Sync

## Why

I can keep writing when signal drops or is lost

## Measurement

users who save locally

-

users who have lost connection

## Requirements

1. Saves offline
   1. GIVEN Writer has no network
      1. WHEN I edit an entry
         1. THEN it is saved locally
2. Keeps uploaded edits
   1. GIVEN Writer has edits saved while offline
      1. WHEN the network returns
         1. THEN the edits upload
      2. WHEN another device opens the entry
         1. THEN it shows the uploaded edits
```

## Rules

[Requirement 0004](docs/requirements/0004-requirement-docs-follow-the-template.md) states the doc rules: IDs, acceptance names and their cap, sections, measurement, steps and never reusing an ID. [Requirement 0005](docs/requirements/0005-each-acceptance-has-one-test.md) states ownership: one test named for each acceptance, no test that proves nothing and a case for each step. The first GIVEN of an acceptance starts with an [actor](../fitness-check-architecture-actors/README.md), as the template shows. IDs follow [ADR 0004](https://github.com/may-journal/fitness-runner/blob/main/fitness/adr-0004/README.md).

[Requirement 0006](docs/requirements/0006-changes-come-paired.md) states that a test and its doc change together. [Requirement 0007](docs/requirements/0007-tests-fit-their-language.md) states how each language finds its tests.

The error for a missing test names the one it expects, and deleting an ID deletes its test. Each language module finds its tests and names the test that owns an ID, as its package documents. A comment or string never counts as a test. Run `fitness explain requirements` for the cap key and its default.

## Steps

An acceptance may chain WHENs as steps under its last GIVEN, as item 2 of the template does. [Requirement 0001](docs/requirements/0001-step-chains-stay-within-the-configured-limit.md) holds the limit and its default, [0004](docs/requirements/0004-requirement-docs-follow-the-template.md) the chain rules and [0005](docs/requirements/0005-each-acceptance-has-one-test.md) the case each step needs. [Requirement 0007](docs/requirements/0007-tests-fit-their-language.md) holds what each language needs, so a chain's cases are proven there.

## Gherkin

The [Cucumber parser](https://github.com/cucumber/gherkin) reads each doc as Gherkin. Shared context sits in the GIVENs once, as the example shows. Requirements [0003](docs/requirements/0003-acceptances-read-as-gherkin.md), [0004](docs/requirements/0004-requirement-docs-follow-the-template.md), [0008](docs/requirements/0008-clauses-stay-split.md) and [0009](docs/requirements/0009-steps-say-each-thing-once.md) hold the numbering, keywords, removed lines, split clauses, alike WHENs and repeated THENs.

```markdown
1. Opens to photos
   1. GIVEN Writer has the app installed on a phone
      1. WHEN they open it
         1. THEN it shows their photos
2. Keeps uploaded edits
   1. GIVEN Writer has the app installed on a phone
   2. GIVEN the user is signed in
      1. WHEN they upload an edited photo
         1. THEN the upload shows its edits
      2. WHEN they reopen the app
         1. THEN the edits are still there
         2. THEN the photo keeps its caption
```

## Behavior

[Requirements](docs/requirements)

## Package roots

A monorepo names its package folders once in `.fitnessrc.json`, and each package root is judged like a repo of its own:

```json
{ "packages": ["cmd/*"] }
```

```text
docs/requirements/0001-the-monorepo-builds.md
cmd/fitness-check-gofmt/docs/requirements/0001-go-files-are-formatted.md
cmd/fitness-check-gofmt/requirements/0001_test.go
```

- Rules: [requirement 0002](docs/requirements/0002-packages-own-their-requirements.md) states how packages own their docs and tests.
- Scope: an acceptance in the repo root of a monorepo that names a `fitness-check-` package fails, since that check's own doc holds it.
- Numbers: a number is unique within its package, so each package starts at 0001.
- Tests: a package's requirement tests sit in its own folder and prove only its docs.
- Missing docs: a package root with code and no docs fails. The repo root may leave its behavior to its packages.
- Pairs: a test pairs with a doc of its own package, as does a product change in `requirements-changed`.
- History: one git log reads every package's docs, and each package keeps its own removed IDs.
- No key: without `packages`, the repo root is the only package, as before.

## Swift tests

Swift requirement tests run on Swift Testing with [apple-ui-testing-library](https://github.com/may-journal/apple-ui-testing-library), hosted in the app, natively as Designed for iPad or on macOS. The [Swift module](https://github.com/may-journal/fitness-runner/blob/main/go/pkg/lang/swift/library.go) judges them.

- Imports: a test file fails on XCTest or another UI test framework, such as ViewInspector or SnapshotTesting.
- Calls: a test file fails on an XCTestCase subclass, Xcode UI tests such as XCUIApplication, or simctl.
- Library: a test file fails unless it, or a file in its folder or below, imports `UITestingLibrary`.
- End-to-end: a file that imports `UITestingLibrarySystem` may use XCTest and simctl, and drives the app through its `SystemSession`.

## Pairs

A requirement test and its doc change together, so neither drifts from the other. Changed behavior extends its existing doc, as [ADR 0004](https://github.com/may-journal/fitness-runner/blob/main/fitness/adr-0004/README.md) sets out. A doc whose words stay the same line by line, as after `biome format`, has not changed. Nor has one whose only change is to HTML comments or deleting its Source line.

- Pairing: a test file pairs with each doc whose number its name or its tests carry, one or several.
- Tests: a changed test file fails unless one of its docs changed too. The error names both files and the missing doc.
- Fixes: a changed test file passes without its doc when it failed this check at base, or only its comments changed.
- Base check: this check judges a checkout of the base commit's docs and test files, where a Swift test importing XCTest fails.
- New files: a test file the base lacks still needs its doc.
- Docs: a changed doc fails unless one of its test files changed too. The error names both files and the missing test.
- Helpers: a test file that names and holds no requirement pairs with nothing.
- No exemptions: there is no trailer, flag, or body text that lets an unpaired change through.

The run decides which changes are judged:

- Commits: the `pre-commit` hook judges the staged files, and HEAD is their base.
- Amends: an amended commit is judged on everything it changes since its parent, which is their base.
- Pull requests: CI judges the diff from the merge base with the base branch, so merged base changes are not judged.
- Replays: a local merge, cherry-pick, or revert, or an amended merge, is not judged again. Its pull request diff still is.
- Full scans: a run with nothing staged, or a push to the default branch, judges no pairs.

## Requirement checks

Some requirements are rules over the repo's own files. A check in `fitness/req-NNNN/`, one of the repo's [local checks](https://github.com/may-journal/fitness-runner/blob/main/docs/local-checks.md), may prove such a doc.

- Ownership: such a doc's acceptances need no test. A test named for one still owns it, and two still fail.
- Orphans: a `fitness/req-NNNN/` folder whose number no doc defines fails.
- Pairs: the check's files pair with the doc as a test file does, so a change to one half alone fails.
- Broken scan: a requirement check that passes having judged no files fails the run.

## Example

```
fitness requirements
fitness requirements --fix
```

Run `fitness explain requirements` for what each key does. This sets each key to its default:

```json
{
  "requirements": {
    "maxAcceptances": 8,
    "maxSteps": 7
  }
}
```

