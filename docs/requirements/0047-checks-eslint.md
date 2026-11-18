# 0047 Checks Eslint

## Why

I keep every JavaScript and TypeScript file clean under one shared lint config, so code reads the same everywhere.

## Measurement

files ESLint reports clean
-
files linted

## Requirements

- 0047.9
    - Given Developer has a repo without JavaScript or TypeScript source, even with a formatting-only `package.json`, or a scoped run staging none
        - When fitness runs eslint
            - Then it passes with nothing to judge
- 0047.2
    - Given Developer has a JavaScript project without `eslint` installed, or with one that crashes printing no JSON
        - When fitness runs eslint
            - Then it shows info with an install hint or ESLint's output, never failing
- 0047.3
    - Given Developer has a project with the shared config and clean JavaScript and TypeScript files
        - When fitness runs eslint
            - Then it passes, counting every tracked file it linted
- 0047.4
    - Given Developer has a JavaScript file with object keys out of order
        - When fitness runs eslint
            - Then it fails with `path:line:col - message (sort-keys)`
- 0047.5
    - Given Developer has the repo's own config that only warns on `console` calls
        - When fitness runs eslint on a file calling `console.log`
            - Then it fails on that warning and ignores the shared rules
- 0047.6
    - Given Developer has the repo's own config that ignores a tracked file
        - When fitness runs eslint
            - Then it fails, asking to remove that ignore pattern
- 0047.7
    - Given Developer has a committed file with findings and a clean staged TypeScript file
        - When fitness runs eslint scoped to staged files
            - Then it passes after linting only the staged file
- 0047.11
    - Given Developer has a `biome.json`, `@biomejs/biome` at any version under any spec, no eslint, and more than 20 Biome findings
        - When fitness runs eslint
            - Then Biome 1.8 or later lints, failing with `path:line:col - message (rule)` per finding; older Biome fails naming its version
