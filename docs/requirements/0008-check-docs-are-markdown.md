# 0008 Check Docs Are Markdown

## Why

I edit each check's rule as Markdown that GitHub renders and the Markdown checks lint, not as Go comments.

## Measurement

checks whose rule is a README

-

checks

## Requirements

```gherkin
SCENARIO: 0008.1
  GIVEN Developer has a check whose rule is a doc.go comment with a heading and list and link and code and spell checker line
  WHEN the doc move script runs on its folder
  THEN its README holds them as Markdown
    AND its doc.go is one line that names and embeds the README
    AND the described summary stays

SCENARIO: 0008.2
  GIVEN Developer has a doc.go that holds more than its comment
  WHEN the doc move script runs on its folder
  THEN it fails naming the folder and leaves both files as they were

SCENARIO: 0008.3
  GIVEN Developer has every check moved to Markdown
  WHEN the README test reads each folder
  THEN each doc.go is one line naming its README and embedding it
    AND the repo README alone links the architecture page

SCENARIO: 0008.4
  GIVEN Developer has a README with a parenthesized link URL or a non-prose first block
  WHEN the check describes itself
  THEN the summary has no stray parenthesis
    AND the summary is empty without prose

SCENARIO: 0008.5
  GIVEN Developer has a doc.go comment that opens with a multi-byte letter
  WHEN the doc move script runs on its folder
  THEN the README opens with that letter as a whole capital

SCENARIO: 0008.6
  GIVEN Developer has each check's README
  WHEN the README test reads it against the check's `--describe`
  THEN the README names only config keys its check declares besides check names
    AND it says to run `fitness explain` when it declares keys
```
