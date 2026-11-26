# 0001 Coverage Measures Every File

## Why

I get Vitest coverage that measures every file, with no exclude entry hiding code from it.

## Measurement

repos whose Vitest coverage excludes nothing

-

repos with Vitest

## Requirements

1. Coverage excludes nothing
   1. GIVEN Developer has a repo with or without a Vitest config
      1. WHEN the repo has no `package.json`
         1. THEN it passes with no files judged
      2. WHEN `coverage.exclude` is empty or missing
         1. THEN it passes
      3. WHEN `coverage.exclude` holds any entry
         1. THEN it fails once per quoted entry telling Developer to remove it
         2. THEN it never reports an imported name
      4. WHEN the repo has a `vitest.config.*` file and a `vitest` field in `package.json`
         1. THEN the config file is the one judged
      5. WHEN the repo has no Vitest config
         1. THEN the shared config is judged
         2. THEN one with an exclude fails
         3. THEN the built in copy passes
