# 0001 Commits Disclose AI Use

## Why

I get every commit message in a repo that uses the convention carrying the AI tool and model that helped.

## Measurement

commit messages that carry both trailers

-

commit messages in repos that use the trailers

## Requirements

1. Commits disclose AI use
   1. GIVEN Developer has commits with or without the trailers `AI-Tools:` or `AI-Models:`
      1. WHEN the repo has no trailers adopted
         1. THEN it passes with no files
      2. WHEN a message carries both trailers with a value each
         1. THEN it passes counting one file
      3. WHEN a repo that uses the trailers has a message missing one
         1. THEN it fails with one error per missing trailer as `commit message missing "<key>:" trailer`
      4. WHEN a trailer has an empty value or sits mid sentence
         1. THEN it counts as missing
      5. WHEN the subject starts with `Merge ` or `Revert `
         1. THEN it passes without trailers
2. Proposed messages are judged
   1. GIVEN Developer has a commit message proposed through `--message`
      1. WHEN the proposed message lacks a trailer in a repo that uses them
         1. THEN it fails for the proposed message instead of the last commit
      2. WHEN the proposed message is empty
         1. THEN it fails asking for both trailers
