# 0036 Checks PR Closes Issue

## Why

Every pull request I merge closes the issues it delivers, so no finished work stays open.

## Measurement

PR bodies that close every issue they deliver
-
PR bodies checked

## Requirements

- 0036.1
    - Given Developer has a body closing an issue with any of GitHub's nine closing keywords, as `#NN`, `owner/repo#NN` or its URL
        - When pr-closes-issue checks the body file
            - Then it passes
- 0036.2
    - Given Developer has a body that only references issues, or puts space after the `#`, as `Closes # 12`
        - When pr-closes-issue checks the body file
            - Then it fails saying it has no closing keyword
- 0036.3
    - Given Developer has a body that implements or plans an issue it does not close
        - When pr-closes-issue checks the body file
            - Then it fails naming each such issue
- 0036.4
    - Given Developer has a closing keyword only inside inline or fenced code or an HTML comment
        - When pr-closes-issue checks the body file
            - Then it fails saying it has no closing keyword
- 0036.5
    - Given Developer has a required issue the body does not close
        - When pr-closes-issue checks it with `--require-close`
            - Then it fails naming that issue
- 0036.6
    - Given Developer has every required issue closed by the body
        - When pr-closes-issue checks it with `--require-close`
            - Then it passes
- 0036.7
    - Given Developer has no body at all
        - When pr-closes-issue runs
            - Then it passes having checked no files
