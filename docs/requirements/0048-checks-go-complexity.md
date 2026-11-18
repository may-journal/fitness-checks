# 0048 Checks Go Complexity

## Why

I keep every Go function small enough to read and test, the same way the house lint rule gates JavaScript.

## Measurement

Go functions within the complexity ceiling
-
Go functions scored

## Requirements

- 0048.1
    - Given Developer has a Go function with more than 5 branches
        - When go-complexity checks the repo
            - Then it fails naming the file, line, function, score, and ceiling
- 0048.2
    - Given Developer has a function with `switch` and `select` cases, `default`, `for`, `range`, `&&`, and `||`
        - When go-complexity scores it
            - Then each adds one point except `default`
- 0048.3
    - Given Developer has a function literal with many branches inside a simple function
        - When go-complexity checks the repo
            - Then it fails naming only the function literal
- 0048.4
    - Given Developer has a method with too many branches
        - When go-complexity checks the repo
            - Then it fails naming the method
- 0048.5
    - Given Developer has a test file with too many branches
        - When go-complexity checks the repo
            - Then it fails like any other Go file
- 0048.6
    - Given Developer has a repo that sets `goComplexity.max` in `.fitnessrc.json`
        - When go-complexity checks the repo
            - Then it judges against that ceiling
- 0048.7
    - Given Developer has a Go file that does not parse
        - When go-complexity checks the repo
            - Then it fails with the parser's error
- 0048.8
    - Given Developer has a repo with no Go files
        - When fitness runs go-complexity
            - Then it passes with 0 files
