# 0001 Go Files Are Formatted

## Why

I get Go files that gofmt would change named, in any repo.

## Measurement

repos whose Go files are gofmt-formatted

-

repos with Go files

## Requirements

1. Unformatted Go files are named with their fix
   1. GIVEN Developer has a repo with Go files
   2. GIVEN the repo may have none
      1. WHEN a Go file is not formatted
         1. THEN it fails as `path: not gofmt-formatted (run: gofmt -w path)`
      2. WHEN every Go file is formatted
         1. THEN it passes
         2. THEN it counts the files checked
      3. WHEN Go files sit where gofmt is not asked to look
         1. THEN a file in `testdata` is not judged
         2. THEN a file in `vendor` is not judged
         3. THEN an untracked file is not judged
      4. WHEN the repo has no Go files
         1. THEN it passes with no files
      5. WHEN a Go file does not parse
         1. THEN it fails as `path:line:column: message`
         2. THEN it never reads as an unformatted file
      6. WHEN Go is not on PATH
         1. THEN it reports info with the hint `see https://go.dev/dl`
      7. WHEN gofmt exits non-zero with no finding
         1. THEN it reports info quoting what went wrong
         2. THEN a crash never reads as an unformatted file
