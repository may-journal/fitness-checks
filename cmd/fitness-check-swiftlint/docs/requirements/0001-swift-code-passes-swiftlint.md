# 0001 Swift Code Passes SwiftLint

## Why

I get SwiftLint judging my Swift code with my repo's own config, in any repo.

## Measurement

repos whose Swift files pass SwiftLint

-

repos with Swift files

## Requirements

- 0001.1
  - Given Developer has a repo, with or without Swift files and a SwiftLint config
    1. When SwiftLint finds a violation
       - Then it fails naming each as `file:line[:col] - reason (rule_id)`
    2. When every Swift file is clean
       - Then it passes, counting the files SwiftLint was asked to lint
    3. When the repo has its own `.swiftlint.yml`
       - Then that config judges, and one excluding every file passes
    4. When a pull request changes some Swift files
       - Then SwiftLint lints only those, unless the config lists what to include
    5. When the repo has no Swift files
       - Then it passes with no files, without running SwiftLint
    6. When SwiftLint is not on PATH
       - Then it reports info with the hint `brew install swiftlint`
    7. When SwiftLint cannot start, is killed, or exits non-zero without a finding
       - Then it reports info quoting what went wrong
