# 0006 Code Is Tested By Its Own Tools

## Why

I get each language's own formatter, linter and test runner judging my code with my repo's config. Every requirement has one test, and a doc and its test change together.

## Measurement

repos whose docs, tests and tool runs all pass their checks

-

repos with code or requirement docs

## Requirements

- 0006.1
  - Given Developer has requirement docs that break the template, step chain or acceptance cap
    - When the requirement doc checks run
      - Then each fails naming the doc and broken rule
- 0006.2
  - Given Developer has requirements or ADRs whose proof is missing, doubled or misnamed
    - When the requirement ownership checks run on the repo
      - Then each one fails naming the requirement and the test or check to add or remove
- 0006.3
  - Given Developer has an unpaired change or a reused ID
    - When the paired change and ID history checks judge my commit or pull request
      - Then each fails naming the files or ID, and a pair or whole restructure passes
- 0006.4
  - Given Developer has Go files with a format, vet or branch-count problem
    - When the Go tooling checks run on them
      - Then each problem fails naming the file and line, and clean code passes
- 0006.5
  - Given Developer has JavaScript or TypeScript files with a lint or format problem, or a disable directive
    - When the Node tooling checks run on them
      - Then each problem fails naming the file and line, and clean code passes
- 0006.6
  - Given Developer has Swift or Java code with a lint or format problem, or a build missing its plugin
    - When the Swift and Java tooling checks run on it
      - Then each problem fails naming the file and the fix, and clean code passes
- 0006.7
  - Given Developer has tests that fail, cover less than all code, or need a runtime version the repo pins
    - When the test and coverage checks run on the pinned runtime
      - Then each shortfall fails naming its cause
- 0006.8
  - Given Developer has Swift pages or views whose names, repeated IDs or localized text break the rules
    - When the Swift Storybook and localized text checks run on the app
      - Then each break fails naming file, line and fix, or passes without Storybook
