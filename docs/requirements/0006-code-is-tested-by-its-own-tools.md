# 0006 Code Is Tested By Its Own Tools

## Why

I get each language's own formatter, linter and test runner judging my code with my repo's config. Every requirement has one test, and a doc and its test change together.

## Measurement

repos whose docs, tests and tool runs all pass their checks

-

repos with code or requirement docs

## Requirements

```gherkin
SCENARIO: 0006.1
  GIVEN Developer has bullet or Gherkin docs breaking the template or step chain or acceptance cap
  WHEN the requirement doc checks run
  THEN each fails naming the doc and rule

SCENARIO: 0006.2
  GIVEN Developer has requirements or ADRs with missing or doubled or misnamed proof or steps no case names
  WHEN the requirement ownership checks run on the repo
  THEN each fails naming the requirement or step or case and what to fix

SCENARIO: 0006.3
  GIVEN Developer has an unpaired change or a reused ID
  WHEN the paired change and ID history checks judge my commit or pull request
  THEN each fails naming the files or ID
    AND a pair or whole restructure passes

SCENARIO: 0006.4
  GIVEN Developer has Go files with a format or vet or branch-count problem
  WHEN the Go tooling checks run on them
  THEN each problem fails naming the file and line
    AND clean code passes

SCENARIO: 0006.5
  GIVEN Developer has JavaScript or TypeScript files with a lint or format problem or a disable directive
  WHEN the Node tooling checks run on them
  THEN each problem fails naming the file and line
    AND clean code passes

SCENARIO: 0006.6
  GIVEN Developer has Java code with a format or style problem or a build missing its plugin or config
  WHEN the Java tooling checks run on it
  THEN each problem fails naming the file and the fix
    AND clean code passes

SCENARIO: 0006.7
  GIVEN Developer has tests that fail or cover less than all code or need a runtime version the repo pins
  WHEN the test and coverage checks run on the pinned runtime
  THEN each shortfall fails naming its cause

SCENARIO: 0006.8
  GIVEN Developer has Swift pages or views whose names or repeated IDs or localized text break the rules
  WHEN the Swift Storybook and localized text checks run on the app
  THEN each break fails naming file and line and fix or passes without Storybook
```
