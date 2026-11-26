# 0001 Markdown Names Follow Camel Case

## Why

I find every doc by a file name in the one case my repo uses.

## Measurement

Markdown files whose names match the repo's case

-

Markdown files in repos that use camelCase

## Requirements

1. Names outside camelCase fail naming the case
   1. GIVEN Developer has a repo whose Markdown names are mostly camelCase
      1. WHEN a name is hyphenated
         1. THEN it fails as `path: filename must be camelCase`
      2. WHEN a name is snake case or PascalCase
         1. THEN each fails as `path: filename must be camelCase`
      3. WHEN every name is camelCase
         1. THEN it passes
      4. WHEN a camelCase doc sits in a capitalized folder
         1. THEN it passes since only the basename is judged
      5. WHEN a commit stages one hyphenated doc
         1. THEN it fails that doc alone
      6. WHEN a name is a standard root doc
         1. THEN it passes
         2. THEN a GitHub template name passes
2. A repo with another case is not judged
   1. GIVEN Developer has a repo whose Markdown names are mostly hyphenated
      1. WHEN most names are hyphenated
         1. THEN it passes with no files
      2. WHEN names tie between hyphens or capitals
         1. THEN kebab-case wins
         2. THEN it passes with no files
