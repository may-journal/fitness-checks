# 0052 Checks Gofmt

## Why

I never review whitespace noise in Go diffs, because every Go file I commit is already formatted the standard way.

## Measurement

Go files `gofmt` leaves unchanged
-
Go files checked

## Requirements

- 0052.1
    - Given Developer has a committed Go file `gofmt` would rewrite
        - When fitness runs `gofmt`
            - Then it fails naming the file and the `gofmt -w` fix
- 0052.2
    - Given Developer has only formatted Go files
        - When fitness runs `gofmt`
            - Then it passes
- 0052.3
    - Given Developer has unformatted Go files under `testdata` or `vendor`
        - When fitness runs `gofmt`
            - Then it passes, since fixtures are not judged
- 0052.4
    - Given Developer has a repo with no Go files
        - When fitness runs `gofmt`
            - Then it passes having checked 0 files
- 0052.5
    - Given Developer has Go files and no Go toolchain on the path
        - When fitness runs `gofmt`
            - Then it shows info with a hint to install Go, never failing
- 0052.6
    - Given Developer has an unformatted Go file that is not tracked by git
        - When fitness runs `gofmt`
            - Then it passes, since only tracked files are judged
- 0052.7
    - Given Developer has a `gofmt` that crashes without listing a file
        - When fitness runs `gofmt`
            - Then it shows info quoting what `gofmt` printed, never failing
- 0052.8
    - Given Developer has a committed Go file `gofmt` cannot parse
        - When fitness runs `gofmt`
            - Then it fails naming the file, line and syntax error
