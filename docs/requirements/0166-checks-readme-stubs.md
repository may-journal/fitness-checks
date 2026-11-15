# 0166 Checks README Stubs

## Why

Each check's README is a short stub that links its full docs, so no generated prose is checked in and none can drift.

## Measurement

Check READMEs that match the stub with links that resolve
-
Check READMEs

## Requirements

- 0166.1
    - Given Developer reads any check folder's README
        - When the stub test reads it
            - Then it holds the check's name, its `--describe` summary, and links to the components view, its `doc.go` and pkg.go.dev, and nothing else
- 0166.2
    - Given Developer follows a stub's relative links
        - When the stub test resolves them from the check folder
            - Then each names a file in this repo
