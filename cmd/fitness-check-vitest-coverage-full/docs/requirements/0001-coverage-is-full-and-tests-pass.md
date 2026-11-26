# 0001 Coverage Is Full And Tests Pass

## Why

I get every JavaScript or TypeScript file covered by tests, because Vitest's thresholds are 100 and its run must pass.

## Measurement

repos whose Vitest run meets full coverage

-

repos with JavaScript or TypeScript source

## Requirements

1. Thresholds are full
   1. GIVEN Developer has a repo with or without JavaScript or TypeScript source
   2. GIVEN it has a project and a shared Vitest config
      1. WHEN the repo has no tracked JavaScript or TypeScript source
         1. THEN it passes with no files
      2. WHEN the project sets a coverage threshold below 100
         1. THEN it fails naming the four metrics
      3. WHEN the shared config sets a coverage threshold below 100
         1. THEN it fails naming the fitness-runner package
      4. WHEN both configs hold 100 and the tests cover all code
         1. THEN it passes
         2. THEN Vitest runs with the repo's own config
2. Vitest runs and reports
   1. GIVEN Developer has a project with full thresholds
      1. WHEN code is untested
         1. THEN it fails showing the coverage line Vitest printed
      2. WHEN Vitest is not installed
         1. THEN it reports info with an install hint
      3. WHEN the binary under `node_modules` cannot run but Vitest is on PATH
         1. THEN the one on PATH runs
      4. WHEN Vitest exits with no finding
         1. THEN it reports info quoting its last lines
