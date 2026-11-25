# 0001 Build Output Stays Out Of Git

## Why

I review and import only source, never the files a build writes.

## Measurement

repos that track no build output and import none

-

repos checked

## Requirements

- 0001.1
  - Given Developer has a repo, with or without `dist`, Maven or Gradle output and TypeScript imports
    1. When git does not ignore the `dist` folder
       - Then it fails with the hint to add `dist/` to `.gitignore`
    2. When git tracks files under any `dist` folder
       - Then one error lists them with the hint `git rm --cached`
    3. When git tracks the build output a language module names, such as Maven `target`
       - Then those files join the same list, and a build folder without Java passes
    4. When a tracked TypeScript file, tests included, imports from a `dist/` path
       - Then it fails as `path:line` for each `from`, `import()` or `require()`
    5. When that import is commented out, or the file is untracked
       - Then it passes
    6. When `dist` is ignored and every import is local
       - Then it passes, counting the scanned sources plus one for the `dist` rules
    7. When the repo has no build output and no TypeScript
       - Then it passes with no files
