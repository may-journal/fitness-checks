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
      - Then its README holds them as Markdown, its doc.go embeds the README, and the described summary stays
- 0008.2
  - Given Developer has a doc.go that holds more than its comment
    - When the doc move script runs on its folder
      - Then it fails naming the folder and leaves both files as they were
- 0008.3
  - Given Developer has the 69 checks moved to Markdown so far
    - When the README test reads each folder
      - Then each doc.go holds a short comment and embeds the README
