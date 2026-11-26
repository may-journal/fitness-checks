# 0001 JavaScript And TypeScript Pass ESLint

## Why

I get ESLint judging my JavaScript and TypeScript with the shared config, or Biome when my repo chose it.

## Measurement

repos whose JavaScript and TypeScript files pass the linter

-

repos with JavaScript or TypeScript files

## Requirements

```gherkin
SCENARIO: 0001.1
  GIVEN Developer has a repo with or without JavaScript or TypeScript files and a lint config
  WHEN a tracked file breaks a rule
  THEN it fails naming each as `path:line:col - message (ruleId)`
  WHEN every tracked file is clean
  THEN it passes and counts the files ESLint was asked to lint
  WHEN the repo has its own `eslint.config.mjs`
  THEN that config judges
  WHEN the repo has no tracked JavaScript or TypeScript source
  THEN it passes with no files and never runs ESLint
  WHEN a pull request changes some files
  THEN ESLint lints only the changed ones and passes with no files when none is lintable
  WHEN ESLint ignores a tracked file
  THEN it fails asking to remove the ignore pattern so the file is linted

SCENARIO: 0001.2
  GIVEN Developer has a repo where ESLint is missing or crashes or where Biome is chosen
  WHEN no eslint binary resolves or ESLint crashes without a finding
  THEN it reports info naming the cause
  WHEN the repo has a `biome.json` and lists `@biomejs/biome` in its dependencies
  THEN Biome lints the same paths with no eslint needed and each error or warning fails as `path:line:col - message (rule)`
  WHEN the Biome found is older than 1.8
  THEN it fails naming its version without running Biome
```
