# 0001 Markdown Names Follow Kebab Case

## Why

I find every doc by a file name in the one case my repo uses.

## Measurement

Markdown files whose names match the repo's case

-

Markdown files in repos that use kebab-case

## Requirements

1. Names outside kebab-case fail naming the case
   1. GIVEN Developer has a repo whose Markdown names are mostly kebab-case
      1. WHEN a name is camelCase
         1. THEN it fails as `path: filename must be kebab-case`
      2. WHEN a name is snake case or PascalCase
         1. THEN each fails as `path: filename must be kebab-case`
         2. THEN a name with a space fails alike
      3. WHEN every name is kebab-case or a standard root doc
         1. THEN it passes
         2. THEN digits in a name pass
         3. THEN a GitHub template name passes
      4. WHEN a kebab-case doc sits in a capitalized folder
         1. THEN it passes since only the basename is judged
      5. WHEN a commit stages one camelCase doc
         1. THEN it fails that doc alone
2. A repo with another case is judged by a tie
   1. GIVEN Developer has a repo whose Markdown names are mostly camelCase
      1. WHEN most names are camelCase
         1. THEN it passes with no files
      2. WHEN names tie between hyphens or capitals
         1. THEN it fails the camelCase file
