# 0001 Build Output Stays Out Of Git

## Why

I review and import only source, never the files a build writes.

## Measurement

repos that track no build output and import none

-

repos checked

## Requirements

```gherkin
SCENARIO: 0001.1
  GIVEN Developer has a repo with or without `dist` or Maven or Gradle output AND TypeScript imports
  WHEN git does not ignore the `dist` folder
  THEN it fails with the hint to add `dist/` to `.gitignore`
  WHEN git tracks files under any `dist` folder
  THEN one error lists them
    AND it gives the hint `git rm --cached`
  WHEN git tracks the build output a language module names such as Maven `target`
  THEN those files join the same list
    BUT a build folder without Java passes
  WHEN a tracked TypeScript file imports from a `dist/` path
    AND the file may be a test
  THEN it fails as `path:line` for each `from` or `import()` or `require()`
  WHEN that import is commented out
    AND the file is untracked
  THEN it passes
  WHEN `dist` is ignored
    AND every import is local
  THEN it passes
    AND it counts the scanned sources plus one for the `dist` rules
  WHEN the repo has no build output
    AND the repo has no TypeScript
  THEN it passes with no files
```
