# 0001 JavaScript And TypeScript Pass ESLint

## Why

I get ESLint judging my JavaScript and TypeScript with the shared config, or Biome when my repo chose it.

## Measurement

repos whose JavaScript and TypeScript files pass the linter

-

repos with JavaScript or TypeScript files

## Requirements

- 0001.1
  - Given Developer has a repo with or without JavaScript or TypeScript files and a lint config
    1. When a tracked file breaks a rule
       - Then it fails naming each as `path:line:col - message (ruleId)`
    2. When every tracked file is clean
       - Then it passes and counts the files ESLint was asked to lint
    3. When the repo has its own `eslint.config.mjs`
       - Then that config judges
    4. When the repo has no tracked JavaScript or TypeScript source
       - Then it passes with no files and never runs ESLint
    5. When a pull request changes some files
       - Then ESLint lints only the changed ones and passes with no files when none is lintable
    6. When ESLint ignores a tracked file
       - Then it fails asking to remove the ignore pattern so the file is linted
- 0001.2
  - Given Developer has a repo where ESLint is missing or crashes or where Biome is chosen
    1. When no eslint binary resolves or ESLint crashes without a finding
       - Then it reports info naming the cause
    2. When the repo has a `biome.json` and lists `@biomejs/biome` in its dependencies
       - Then Biome lints the same paths with no eslint needed and each error or warning fails as `path:line:col - message (rule)`
    3. When the Biome found is older than 1.8
       - Then it fails naming its version without running Biome
