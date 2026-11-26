# 0006 Code Is Tested By Its Own Tools

## Why

I get each language's own formatter, linter and test runner judging my code with my repo's config. Every requirement has one test, and a doc and its test change together.

## Measurement

repos whose docs, tests and tool runs all pass their checks

-

repos with code or requirement docs

## Requirements

1. Removed in 2026.11.26.0200: moved to the requirements doc 0004
2. Each requirement has one proof
   1. GIVEN Developer has unproven or doubly proven requirements
   2. GIVEN a proof is misnamed or a step has no case
   3. GIVEN a check judges no files
      1. WHEN the ownership checks run on the repo
         1. THEN each fails naming its requirement or step and fix
3. Changes come paired
   1. GIVEN Developer has an unpaired change
      1. WHEN the history checks judge my commit
         1. THEN each fails naming the files
         2. THEN a pair passes
4. Removed in 2026.11.26.0200: moved to the go-test and go-vet and gofmt and go-complexity docs
5. Removed in 2026.11.26.0300: moved to the no-eslint-disable doc
6. Removed in 2026.11.26.0200: moved to the spotless doc
7. Tests pass and cover
   1. GIVEN Developer has tests that cover less than all code
      1. WHEN the test and coverage checks run
         1. THEN each shortfall fails naming its cause
8. Removed in 2026.11.26.0200: moved to the swift-storybook-requirements doc
