# 0012 Checks README Stubs

## Why

Each check's README is a short stub that links its full docs, so no generated prose is checked in. A test catches a stub whose summary drifts from `doc.go`.

## Measurement

Check READMEs that match the stub with links that resolve
-
Check READMEs

## Requirements

- 0012.1
    - Given Developer reads any check folder's README
        - When the stub test reads it
            - Then it holds the check's name, the summary `--describe` prints, and links to the components view, its `doc.go` and pkg.go.dev, and nothing else
- 0012.2
    - Given Developer follows a stub's relative links
        - When the stub test resolves them from the check folder
            - Then each names a file in this repo
