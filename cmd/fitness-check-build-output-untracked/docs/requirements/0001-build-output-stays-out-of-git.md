# 0001 Build Output Stays Out Of Git

## Why

I review and import only source, never the files a build writes.

## Measurement

repos that track no build output and import none

-

repos checked

## Requirements

1. Build output stays out of git
   1. GIVEN Developer has a repo
   2. GIVEN it may hold `dist` output
   3. GIVEN it may hold Maven or Gradle output
   4. GIVEN it may hold TypeScript imports
      1. WHEN git does not ignore the `dist` folder
         1. THEN it fails with the hint to add `dist/` to `.gitignore`
      2. WHEN git tracks files under any `dist` folder
         1. THEN one error lists them
         2. THEN it hints `git rm --cached`
      3. WHEN git tracks the build output a language module names such as Maven `target`
         1. THEN those files join the same list
         2. THEN a build folder without Java passes
      4. WHEN a tracked TypeScript file imports from a `dist/` path
         1. THEN it fails as `path:line` for each `from`
         2. THEN it fails the same way for each `import()`
         3. THEN it fails the same way for each `require()`
         4. THEN tests are judged like any other file
      5. WHEN no tracked file imports from a `dist/` path
         1. THEN a commented-out import passes
         2. THEN an untracked file passes too
         3. THEN an ignored `dist` with every import local passes
         4. THEN it counts the scanned sources plus one for the `dist` rules
      6. WHEN the repo has neither build output nor TypeScript
         1. THEN it passes with no files
