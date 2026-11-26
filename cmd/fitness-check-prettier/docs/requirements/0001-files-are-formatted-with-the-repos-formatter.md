# 0001 Files Are Formatted With The Repo's Formatter

## Why

I see every tracked file formatted by my repo's own formatter, Prettier or Biome, and told when it is not.

## Measurement

repos whose tracked files are formatted

-

repos with a formatter

## Requirements

1. Files pass Prettier with the repo's config
   1. GIVEN Developer has a repo that may hold a `package.json`
   2. GIVEN it may hold a Prettier config
      1. WHEN the repo has its own Prettier config
         1. THEN `prettier` judges every tracked file with it
         2. THEN a file it would reformat fails by name
      2. WHEN the repo has no `package.json`
         1. THEN it passes without running Prettier
      3. WHEN the repo has no Prettier config
         1. THEN the shared config judges with single quotes
         2. THEN the synced `.github` templates keep Prettier defaults
      4. WHEN files are staged
         1. THEN only the staged files are checked
         2. THEN a symbolic link is left out
      5. WHEN a `--write` argument is passed
         1. THEN Prettier writes in place instead of checking
2. Prettier trouble is reported
   1. GIVEN Developer has a repo that Prettier cannot judge cleanly
      1. WHEN a `.prettierignore` exists
         1. THEN it fails telling to delete it and fix the findings
      2. WHEN Prettier is not installed
         1. THEN it reports info with an install hint
      3. WHEN a plugin is missing
         1. THEN it reports info quoting Prettier's error
      4. WHEN a file cannot be parsed
         1. THEN it fails naming the file
3. Biome formats a repo that chose it
   1. GIVEN Developer has a repo with `biome.json` and `@biomejs/biome` in its dependencies
      1. WHEN a file needs reformatting
         1. THEN `prettier` fails naming the file
         2. THEN it needs no Prettier
      2. WHEN Biome is older than 1.8
         1. THEN it fails naming the version
         2. THEN it never runs Biome
      3. WHEN a `--write` argument is passed
         1. THEN Biome writes in place
         2. THEN it gets only its own flags
      4. WHEN a Prettier flag has no Biome counterpart
         1. THEN it fails naming the flag
