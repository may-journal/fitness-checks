# 0008 Check Docs Are Markdown

## Why

I edit each check's rule as Markdown that GitHub renders and the Markdown checks lint, not as Go comments.

## Measurement

checks whose rule is a README

-

checks

## Requirements

- 0008.1
  - Given Developer has a check whose rule is a doc.go comment with a heading, list, link, code and spell checker line
    - When the doc move script runs on its folder
      - Then its README holds them as Markdown, its doc.go is one line that names and embeds the README, and the described summary stays
- 0008.2
  - Given Developer has a doc.go that holds more than its comment
    - When the doc move script runs on its folder
      - Then it fails naming the folder and leaves both files as they were
- 0008.3
  - Given Developer has the checks moved to Markdown so far
    - When the README test reads each folder
      - Then each doc.go is one line naming its README and embedding it, and the repo README alone links the architecture page
- 0008.4
  - Given Developer has a README with a parenthesized link URL, or a non-prose first block
    - When the check describes itself
      - Then the summary has no stray parenthesis, and is empty without prose
- 0008.5
  - Given Developer has a doc.go comment that opens with a multi-byte letter
    - When the doc move script runs on its folder
      - Then the README opens with that letter as a whole capital
- 0008.6
  - Given Developer has each check's README
    - When the README test reads it against the check's `--describe`
      - Then the README names no config key that the check does not declare, and says to run `fitness explain` when it declares keys
