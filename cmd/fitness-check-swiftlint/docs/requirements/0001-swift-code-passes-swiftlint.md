# 0001 Swift Code Passes SwiftLint

## Why

I get SwiftLint judging my Swift code with my repo's own config, in any repo.

## Measurement

repos whose Swift files pass SwiftLint

-

repos with Swift files

## Requirements

```gherkin
SCENARIO: 0001.1
  GIVEN Developer has a repo with or without Swift files and a SwiftLint config
  WHEN SwiftLint finds a violation
  THEN it fails naming each as `file:line[:col] - reason (rule_id)`
  WHEN every Swift file is clean
  THEN it passes and counts the files SwiftLint was asked to lint
  WHEN the repo has its own `.swiftlint.yml`
  THEN that config judges and one excluding every file passes
  WHEN a pull request changes some Swift files
  THEN SwiftLint lints only those unless the config lists what to include
  WHEN the repo has no Swift files
  THEN it passes with no files and never runs SwiftLint
  WHEN SwiftLint is not on PATH
  THEN it reports info with the hint `brew install swiftlint`
  WHEN SwiftLint cannot start or is killed or exits non-zero without a finding
  THEN it reports info quoting what went wrong
```
