# 0001 Swift Code Passes SwiftLint

## Why

I get SwiftLint judging my Swift code with my repo's own config, in any repo.

## Measurement

repos whose Swift files pass SwiftLint

-

repos with Swift files

## Requirements

1. Violations fail
   1. GIVEN Developer has a repo with or without Swift files
   2. GIVEN the repo has a SwiftLint config
      1. WHEN SwiftLint finds a violation
         1. THEN it fails naming each as `file:line[:col] - reason (rule_id)`
      2. WHEN every Swift file is clean
         1. THEN it passes
         2. THEN it counts the files SwiftLint was asked to lint
      3. WHEN the repo has its own `.swiftlint.yml`
         1. THEN that config judges
         2. THEN a config excluding every file passes
      4. WHEN a pull request changes some Swift files
         1. THEN SwiftLint lints only those unless the config lists what to include
      5. WHEN the repo has no Swift files
         1. THEN it passes with no files
         2. THEN it never runs SwiftLint
      6. WHEN SwiftLint is not on PATH
         1. THEN it reports info with the hint `brew install swiftlint`
      7. WHEN SwiftLint cannot start or is killed or exits nonzero without a finding
         1. THEN it reports info quoting what went wrong
