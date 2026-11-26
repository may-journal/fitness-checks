# 0001 Node Matches The Repo Pin

## Why

I am told when my active Node is older than the version my repo pins.

## Measurement

runs where the active Node satisfies `.nvmrc`

-

runs of `node-version` in repos with a `package.json`

## Requirements

1. A repo's Node meets its pin
   1. GIVEN Developer has a repo that may hold a `package.json`
   2. GIVEN it may pin a Node version in `.nvmrc`
      1. WHEN the repo has no `package.json`
         1. THEN `node-version` passes with zero files
      2. WHEN the repo has no `.nvmrc`
         1. THEN it fails with `missing .nvmrc`
      3. WHEN `.nvmrc` holds a major no higher than the active Node's
         1. THEN it passes
         2. THEN it reads a `v` prefix as no prefix
      4. WHEN `.nvmrc` holds a major above the active Node's
         1. THEN it fails with `does not satisfy .nvmrc`
         2. THEN it gives the fix `nvm use`
      5. WHEN `.nvmrc` holds no number
         1. THEN it fails as an unsatisfiable `NaN.x`
      6. WHEN `node` is not on PATH
         1. THEN it reports info with an install hint
