# 0003 Acceptances Read As Gherkin

## Why

I write each acceptance as Gherkin with keywords in caps and one condition per line, and a standard parser checks its shape.

## Measurement

Gherkin acceptances the check reads and pairs with their tests

-

Gherkin acceptances

## Requirements

```gherkin
SCENARIO: 0003.1
  GIVEN Developer has a requirement doc whose acceptances sit in one all-caps Gherkin block
  WHEN `requirements` reads it
  THEN it passes and pairs each acceptance and step with its test just as for bullets

SCENARIO: 0003.2
  GIVEN Developer has a Gherkin keyword that is not in caps
  WHEN `requirements` reads it
  THEN it fails naming the line and the keyword

SCENARIO: 0003.3
  GIVEN Developer has an AND line no deeper than the step it extends
  WHEN `requirements` reads it
  THEN it fails naming the line

SCENARIO: 0003.4
  GIVEN Developer has a Gherkin step that holds a comma outside code
  WHEN `requirements` reads it
  THEN it fails asking for an AND line

SCENARIO: 0003.5
  GIVEN Developer has a Gherkin block the Cucumber parser rejects
  WHEN `requirements` reads it
  THEN it fails at the doc line the parser names

SCENARIO: 0003.6
  GIVEN Developer has a doc of bullet acceptances
  WHEN migration `0003.6-gherkin` runs up through `--fix` or `--up`
  THEN each acceptance becomes an all-caps SCENARIO in one Gherkin block that passes
  WHEN it runs down through `--down`
  THEN the doc holds bullet acceptances again with each AND joined to its line

SCENARIO: 0003.7
  GIVEN Developer moves a bullet doc to Gherkin and splits its commas into AND lines
  WHEN `requirements` judges the change
  THEN the move pairs as a reformat so its tests need no change
```
