# 0001 Source Holds No Eslint Disable

## Why

I fix the lint rule and never silence it. No JS or TS source holds a disable directive.

## Measurement

JS and TS source files with no ESLint disable directive

-

JS and TS source files

## Requirements

1. Directives fail by line
   1. GIVEN Developer has JS or TS source files
   2. GIVEN some hold ESLint directives
      1. WHEN a file holds an `eslint-disable` directive
         1. THEN it fails once per hit
         2. THEN each hit names the file
         3. THEN each hit names the line
         4. THEN each hit names the directive
         5. THEN `eslint-disable-line` counts as a hit
         6. THEN `eslint-disable-next-line` counts as a hit
      2. WHEN no file holds a directive
         1. THEN it passes
         2. THEN it counts the files scanned
      3. WHEN directive text sits in a string in a file of any scanned extension
         1. THEN it fails like any other hit
      4. WHEN directive text sits in a Markdown file
         1. THEN the file passes
      5. WHEN the repo has no JS or TS source
         1. THEN the check passes with no files
      6. WHEN `node_modules` or `dist` holds a directive
         1. THEN it fails like any other folder
2. Only tracked files are scanned
   1. GIVEN Developer has a repo with staged or untracked files
      1. WHEN a file holding a directive is untracked
         1. THEN it passes
      2. WHEN only one file is staged
         1. THEN only that file is scanned
