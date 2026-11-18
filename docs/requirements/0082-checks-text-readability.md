# 0082 Checks Text Readability

## Why

I catch run-on docs and notation posing as prose, without being punished for ordinary dense writing.

## Measurement

docs with 2 of 3 scores in alarm reported red
-
docs with 2 of 3 scores in alarm

## Requirements

- 0082.1
    - Given Developer has a doc of long words in one endless sentence
        - When text-readability checks the repo
            - Then it fails naming the doc, its scores, and how to fix it
- 0082.2
    - Given Developer has a doc of plain short sentences
        - When text-readability checks the repo
            - Then it passes
- 0082.3
    - Given Developer has a dense doc under 100 prose words
        - When text-readability checks the repo
            - Then it passes, since short docs are not judged
- 0082.4
    - Given Developer has a doc where only one score reaches its alarm
        - When text-readability checks the repo
            - Then it passes, since 2 of 3 must agree
- 0082.5
    - Given Developer has dense notation inside a fenced code block
        - When text-readability checks the repo
            - Then it passes, since code is not scored
- 0082.6
    - Given Developer has a dense issue or PR body
        - When text-readability checks it in body mode
            - Then it fails naming the description
- 0082.7
    - Given Developer has a repo that lowers its own `textReadability` bands
        - When text-readability checks plain prose
            - Then it fails against those bands
- 0082.8
    - Given Developer has the `--report` flag
        - When text-readability checks the repo
            - Then it prints each doc's word count without changing the verdict
