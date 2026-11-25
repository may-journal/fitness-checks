# 0001 Unknown Words Fail

## Why

I get every misspelling and forbidden phrase in my repo named, with project words always winning over the shared ones.

## Measurement

repos whose checked files hold no unknown or forbidden words

-

repos with checked files

## Requirements

```gherkin
SCENARIO: 0001.1
  GIVEN Developer has a repo with tracked files and an optional `cspell.json`
  WHEN a checked file holds an unknown word
  THEN it fails naming each as `path:line:col - Unknown word (word)`
  WHEN a file is over 64 KB
  THEN it is checked in line-aligned chunks with its `cspell:` directives applied to the whole file
  WHEN the repo's own `cspell.json` lists a word or an installed shared package lists another
  THEN the repo's list wins and the owner and agent names pass
  WHEN the resolved `cspell.json` sets `ignorePaths`
  THEN it fails before any scan telling Developer to remove it and fix the words
  WHEN files are binary or marked `linguist-generated` or are lockfiles or build wrapper files
  THEN they are left out and nothing else is
  WHEN a scoped run has staged changes
  THEN only the changed paths that still exist on disk are checked
  WHEN an issue or PR body holds an unknown word
  THEN it fails naming `(description):line:col`
  WHEN the repo has no `cspell.json`
  THEN an installed `@mayjournal/fitness-shared` package or the copy embedded in the binary supplies the words
  WHEN a `flagWords` word or phrase appears whole in any case
  THEN it fails as `path:line:col - Forbidden word (entry)` except in a `cspell.json` file
  WHEN a PR body holds Renovate's hidden debug comment
  THEN that comment is skipped and the rest is still spelled
```
