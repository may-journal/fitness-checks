# 0001 Ignore Patterns Say Why

## Why

I get every `.gitignore` pattern explained by a comment, so ignores stay easy to audit and onboard.

## Measurement

repos whose `.gitignore` patterns each follow an explanatory comment

-

repos with a `.gitignore`

## Requirements

```gherkin
SCENARIO: 0001.1
  GIVEN Developer has a repo with or without a `.gitignore`
  WHEN there is no `.gitignore`
  THEN it passes with no files
  WHEN every pattern has a `#` comment with text on the line directly above it
  THEN it passes and counts the one file
  WHEN a pattern has no comment above it
  THEN it fails naming the line and the pattern as `.gitignore:<line>: pattern "<text>" has no explanatory # comment on the line above`
  WHEN one comment sits above two patterns
  THEN it fails naming the second pattern
  WHEN a blank line sits between the comment and its pattern
  THEN it fails naming the pattern
  WHEN the comment above a pattern is a bare `#`
  THEN it fails naming the pattern
```
