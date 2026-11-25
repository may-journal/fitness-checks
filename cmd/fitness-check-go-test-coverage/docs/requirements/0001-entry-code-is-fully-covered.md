# 0001 Entry Code Is Fully Covered

## Why

I get every command package and declared library entry reaching full statement coverage and a report of tests that overlap or run unstably.

## Measurement

entry packages at full statement coverage

-

entry packages

## Requirements

```gherkin
SCENARIO: 0001.1
  GIVEN Developer has a repo with or without Go modules
  WHEN the repo has no Go modules
  THEN it passes with no files
  WHEN Go is not on PATH
  THEN it reports info with the hint `see https://go.dev/dl`
  WHEN a module has a library and no declared entry
  THEN it fails asking to declare entries
  WHEN `--entry` names an unknown package
  THEN it fails naming that package
  WHEN an entry runs fewer statements than it holds
  THEN it fails naming each entry with its covered and total statements
  WHEN a top-level test passes a named helper to `t.Run`
  THEN it fails naming the file and line and the helper
  WHEN `--audit-map` names a new path outside the repo
  THEN it passes and writes the source map there
  WHEN a test covers different blocks on each run
  THEN it reports the test as unstable and excludes it from overlap advice
  WHEN the audit runs parallel tests and serial tests
  THEN parallel tests run alongside each other and serial packages run side by side but never within one package
```
