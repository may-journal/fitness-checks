# 0001 Go Functions Stay Simple

## Why

I get Go functions held to a cyclomatic complexity ceiling like the house eslint rule so none grows too tangled to read.

## Measurement

Go functions at or under the ceiling

-

Go functions

## Requirements

```gherkin
SCENARIO: 0001.1
  GIVEN Developer has a repo with or without Go files
  WHEN a function has more branches than the ceiling allows
  THEN it fails naming `file:line` and its complexity and the maximum allowed
  WHEN a function holds if and for and range and switch and select cases and `&&` and `||`
  THEN each adds one point to a start of 1 and default cases add none
  WHEN a function literal sits in another function
  THEN it scores alone and the enclosing function does not carry its branches
  WHEN a method has too many branches
  THEN it fails reported as `method <name>`
  WHEN a `_test.go` file has a function over the ceiling
  THEN it fails like any other Go file
  WHEN `.fitnessrc.json` sets `goComplexity.max`
  THEN that ceiling judges instead of 5
  WHEN a Go file does not parse
  THEN it fails with the parser's error
  WHEN the repo has no Go files
  THEN it passes with no files
```
