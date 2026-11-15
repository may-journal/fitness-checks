# 0167 Checks Workflow Scripts

## Why

Shell inside a workflow cannot be linted or tested. So every step that does real work calls a script file.

## Measurement

workflow run steps that are one trivial command or call a script
-
workflow run steps

## Requirements

- 0167.1
    - Given Developer has a workflow and an action whose run steps are one trivial command or call a script
        - When workflow-scripts runs
            - Then it passes
- 0167.2
    - Given Developer has a workflow step with several lines of shell
        - When workflow-scripts runs
            - Then it fails naming the file, the line, and the script to write
- 0167.3
    - Given Developer has an action step with a pipe, a chain or control flow on one line
        - When workflow-scripts runs
            - Then it fails naming that operator
- 0167.4
    - Given Developer has a step whose operators sit only in quotes or a GitHub expression
        - When workflow-scripts runs
            - Then it passes
- 0167.5
    - Given Developer has a step whose shell spans lines as a plain scalar, or whose block holds a `run:` line
        - When workflow-scripts runs
            - Then it fails the step by its lines, and leaves text inside a block alone
