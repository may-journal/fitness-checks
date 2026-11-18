# 0055 Checks Node Version

## Why

I learn at once when my Node is older than the repo's `.nvmrc` asks, before a build fails strangely.

## Measurement

JavaScript repos whose Node satisfies `.nvmrc`
-
JavaScript repos checked

## Requirements

- 0055.1
    - Given Developer has a repo with no `package.json`
        - When fitness runs node-version
            - Then it passes having checked no files
- 0055.2
    - Given Developer has a JavaScript repo without an `.nvmrc`
        - When fitness runs node-version
            - Then it fails with `missing .nvmrc`
- 0055.3
    - Given Developer has an `.nvmrc` major at or below the installed Node
        - When fitness runs node-version
            - Then it passes
- 0055.4
    - Given Developer has an `.nvmrc` with a `v` prefix
        - When fitness runs node-version
            - Then it reads the version the same way
- 0055.5
    - Given Developer has an `.nvmrc` major above the installed Node
        - When fitness runs node-version
            - Then it fails naming the required major and `nvm use`
- 0055.6
    - Given Developer has an `.nvmrc` with no version number
        - When fitness runs node-version
            - Then it fails reporting the requirement as `NaN.x`
- 0055.7
    - Given Developer has no `node` on the `PATH`
        - When fitness runs node-version
            - Then it shows info telling me to install Node, never failing
