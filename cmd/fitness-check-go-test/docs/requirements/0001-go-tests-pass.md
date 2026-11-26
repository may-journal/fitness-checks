# 0001 Go Tests Pass

## Why

I get every Go test in my repo run and each failure named, in any repo.

## Measurement

repos whose Go tests pass

-

repos with Go modules

## Requirements

1. Go tests run and each failure is named
   1. GIVEN Developer has a repo with Go modules
   2. GIVEN the repo may have none
      1. WHEN a test fails
         1. THEN it fails naming the test and its `file.go:line` message
      2. WHEN the failing test sits in a module below the repo root
         1. THEN each line is prefixed with that module's folder
      3. WHEN a test does not compile
         1. THEN it fails with the build error and its `file.go:line:col`
      4. WHEN every test passes
         1. THEN it passes and counts the modules tested
      5. WHEN the repo has no Go module to test
         1. THEN it passes with no files and never runs `go test`
         2. THEN a `testdata` module counts as none
         3. THEN a `vendor` module counts as none
         4. THEN a change touching no Go file counts as none
      6. WHEN Go is not on PATH
         1. THEN it reports info with the hint `see https://go.dev/dl`
      7. WHEN `go test` fails without a finding
         1. THEN it reports info quoting what went wrong
         2. THEN a failure to start and a kill count the same
