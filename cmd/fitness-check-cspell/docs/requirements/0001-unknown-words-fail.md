# 0001 Unknown Words Fail

## Why

I get every misspelling and forbidden phrase in my repo named, with project words always winning over the shared ones.

## Measurement

repos whose checked files hold no unknown or forbidden words

-

repos with checked files

## Requirements

1. Unknown words fail
   1. GIVEN Developer has a repo with tracked files
   2. GIVEN it may hold a `cspell.json`
      1. WHEN a checked file holds an unknown word
         1. THEN it fails naming each as `path:line:col - Unknown word (word)`
      2. WHEN a file is over 64 KB
         1. THEN it is checked in line-aligned chunks
         2. THEN its `cspell:` directives apply to the whole file
2. The repo's word list wins
   1. GIVEN Developer has a repo with tracked files
   2. GIVEN it has a `cspell.json`
      1. WHEN the repo's own list holds a word while an installed shared package holds another
         1. THEN the repo's list wins
         2. THEN the owner names pass
         3. THEN the agent names pass
      2. WHEN the resolved `cspell.json` sets `ignorePaths`
         1. THEN it fails before any scan
         2. THEN it tells Developer to remove it
         3. THEN it tells Developer to fix the words
3. Only the right files are read
   1. GIVEN Developer has a repo with tracked files
      1. WHEN files are binary or marked `linguist-generated` or are lockfiles or build wrapper files
         1. THEN they are left out
         2. THEN nothing else is
      2. WHEN a scoped run has staged changes
         1. THEN only the changed paths that still exist on disk are checked
4. Fallback words and forbidden phrases apply
   1. GIVEN Developer has a repo with tracked files
      1. WHEN the repo has no `cspell.json`
         1. THEN an installed `@mayjournal/fitness-shared` package supplies the words
         2. THEN the copy embedded in the binary supplies them if none is installed
      2. WHEN a `flagWords` word or phrase appears whole in any case
         1. THEN it fails as `path:line:col - Forbidden word (entry)`
         2. THEN a `cspell.json` file is exempt
5. Bodies are spelled
   1. GIVEN Developer has an issue or PR body
      1. WHEN the body holds an unknown word
         1. THEN it fails naming `(description):line:col`
      2. WHEN a PR body holds Renovate's hidden debug comment
         1. THEN that comment is skipped
         2. THEN the rest is still spelled
