# 0001 Go Functions Stay Simple

## Why

I get Go functions held to a cyclomatic complexity ceiling like the house eslint rule so none grows too tangled to read.

## Measurement

Go functions at or under the ceiling

-

Go functions

## Requirements

1. Functions over the ceiling fail
   1. GIVEN Developer has a repo with Go files
      1. WHEN a function has more branches than the ceiling allows
         1. THEN it fails naming `file:line`
         2. THEN it names its complexity and the maximum allowed
      2. WHEN a function holds `if` and `for` and `range` and `switch` and `select` cases and `&&` and `||`
         1. THEN each adds one point to a start of 1
         2. THEN default cases add none
      3. WHEN a function literal sits in another function
         1. THEN it scores alone
         2. THEN the enclosing function does not carry its branches
      4. WHEN a method has too many branches
         1. THEN it fails reported as `method <name>`
2. Every Go file is judged by the repo's ceiling
   1. GIVEN Developer has a repo with or without Go files
      1. WHEN a `_test.go` file has a function over the ceiling
         1. THEN it fails like any other Go file
      2. WHEN `.fitnessrc.json` sets `goComplexity.max`
         1. THEN that ceiling judges instead of 5
      3. WHEN a Go file does not parse
         1. THEN it fails with the parser's error
      4. WHEN the repo has no Go files
         1. THEN it passes with no files
