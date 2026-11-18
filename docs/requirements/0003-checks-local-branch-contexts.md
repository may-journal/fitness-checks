# 0003 Local Branch Contexts

## Why

My push checks its PR and Plans before Actions. Fitness holds the judgments, so this program only declares the check.

## Measurement

local-branch answers that report what they could judge
-
local-branch answers

## Requirements

- 0003.1
    - Given Developer has the local-branch program built from this repo
        - When I run it with `--describe`, and run it alone
            - Then it names the check with protocol 1 and declares its scope, and the run exits 3 saying it runs only inside fitness
