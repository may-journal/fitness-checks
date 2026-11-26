# 0001 JavaScript And TypeScript Pass ESLint

## Why

I get ESLint judging my JavaScript and TypeScript with the shared config, or Biome when my repo chose it.

## Measurement

repos whose JavaScript and TypeScript files pass the linter

-

repos with JavaScript or TypeScript files

## Requirements

1. Rule breaks fail
   1. GIVEN Developer has a repo with JavaScript or TypeScript files and a lint config
      1. WHEN every tracked file is clean
         1. THEN it passes
         2. THEN it counts the files ESLint was asked to lint
      2. WHEN a tracked file breaks a rule
         1. THEN it fails naming each as `path:line:col - message (ruleId)`
      3. WHEN the repo has its own `eslint.config.mjs`
         1. THEN that config judges in place of the shared rules
      4. WHEN ESLint ignores a tracked file
         1. THEN it fails asking to remove the ignore pattern so the file is linted
2. Only lintable files are linted
   1. GIVEN Developer has a repo with or without JavaScript or TypeScript files
      1. WHEN the repo has no tracked JavaScript or TypeScript source
         1. THEN it passes with no files
         2. THEN it never runs ESLint
      2. WHEN a pull request or commit changes some files
         1. THEN ESLint lints only the changed ones
3. A missing ESLint informs and Biome stands in
   1. GIVEN Developer has a repo with JavaScript files where ESLint is missing or crashes or Biome is chosen
      1. WHEN no eslint binary resolves or ESLint crashes without a finding
         1. THEN it reports info naming the cause
      2. WHEN the repo has a `biome.json` and lists `@biomejs/biome` in its dependencies
         1. THEN Biome lints the same paths with no eslint needed
         2. THEN each error or warning fails as `path:line:col - message (rule)`
         3. THEN a Biome older than 1.8 fails naming its version without running Biome
