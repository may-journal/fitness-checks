# 0001 Changelog Bullets Are Few, Short And Typed

## Why

I read each change's entry at a glance: a few short bullets, each naming its type.

## Measurement

entry files whose sections pass every bullet rule

-

entry files

## Requirements

1. Entry sections stay short and typed
   1. GIVEN Developer has a repo
   2. GIVEN it may hold `changelog/` entry files
      1. WHEN each section holds 3 to 5 bullets under 365 characters each opening with a semantic type
         1. THEN it passes counting the entry files
         2. THEN types read `Feat:` or `feat(scope)!:` alike
      2. WHEN a `###` section holds fewer than 3 `- ` bullets
         1. THEN it fails naming the file
         2. THEN it fails naming the section
         3. THEN it fails naming its count
         4. THEN a section of more than 5 fails the same way
      3. WHEN a bullet reaches 365 characters with its continuation lines joined
         1. THEN it fails naming the file
         2. THEN it fails naming the line
         3. THEN it fails naming the length
      4. WHEN a bullet opens with none of the types `Feat:` through `Revert:`
         1. THEN it fails that bullet alone naming its line
         2. THEN a commit subject of such a type passes
      5. WHEN several entry files break these rules
         1. THEN each error names its own file
      6. WHEN the repo has only a `CHANGELOG.md`
         1. THEN it passes with no files
