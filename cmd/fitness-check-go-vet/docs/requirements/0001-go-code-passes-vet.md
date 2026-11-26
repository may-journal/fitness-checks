# 0001 Go Code Passes Vet

## Why

I get suspicious Go code the compiler accepts caught, in any repo.

## Measurement

repos whose Go modules pass `go vet`

-

repos with Go modules

## Requirements

1. Go vet problems are named
   1. GIVEN Developer has a repo with Go modules
   2. GIVEN the repo may have none
      1. WHEN `go vet` finds a problem in a module at the repo root
         1. THEN it fails naming the finding as `file.go:line:col: message`
      2. WHEN the module sits below the repo root
         1. THEN each finding is prefixed with that module's folder
      3. WHEN every module is clean
         1. THEN it passes
         2. THEN it counts the modules vetted
      4. WHEN the repo has no Go module to vet
         1. THEN it passes with no files
         2. THEN it never runs `go vet`
         3. THEN a `testdata` module counts as none
         4. THEN a `vendor` module counts as none
      5. WHEN a run is scoped to changed files with none of them Go
         1. THEN it passes with no files
      6. WHEN Go is not on PATH
         1. THEN it reports info with the hint `see https://go.dev/dl`
      7. WHEN `go vet` fails naming no `.go` file
         1. THEN it reports info quoting what went wrong
         2. THEN a failure to start counts the same
         3. THEN a kill counts the same
