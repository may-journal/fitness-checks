# 0001 Commits Disclose AI Use

## Why

I get every commit message in a repo that uses the convention carrying the AI tool and model that helped.

## Measurement

commit messages that carry both trailers

-

commit messages in repos that use the trailers

## Requirements

```gherkin
SCENARIO: 0001.1
  GIVEN Developer has commits with or without the `AI-Tools:` and `AI-Models:` trailers
  WHEN neither the message nor the last 100 commits carry a trailer or there is no repo
  THEN it passes with no files
  WHEN a message carries both trailers with a value each
  THEN it passes counting one file
  WHEN a repo that uses the trailers has a message missing one
  THEN it fails with one error per missing trailer as `commit message missing "<key>:" trailer`
  WHEN a trailer has an empty value or sits mid sentence
  THEN it counts as missing
  WHEN the subject starts with `Merge ` or `Revert `
  THEN it passes without trailers

SCENARIO: 0001.2
  GIVEN Developer has a commit message proposed through `--message`
  WHEN the proposed message lacks a trailer in a repo that uses them
  THEN it fails for the proposed message instead of the last commit
  WHEN the proposed message is empty
  THEN it fails asking for both trailers
```
