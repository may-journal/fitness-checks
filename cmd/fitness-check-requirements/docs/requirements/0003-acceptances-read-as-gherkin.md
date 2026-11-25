# 0003 Acceptances Read As Gherkin

## Why

I write each acceptance as Gherkin with keywords in caps and one condition per line, and a standard parser checks its shape.

## Measurement

Gherkin acceptances the check reads and pairs with their tests

-

Gherkin acceptances

## Requirements

- 0003.1
  - Given Developer has a requirement doc whose acceptances sit in one all-caps Gherkin block
    - When `requirements` reads it
      - Then it passes and pairs each acceptance and step with its test just as for bullets
- 0003.2
  - Given Developer has a Gherkin keyword that is not in caps
    - When `requirements` reads it
      - Then it fails naming the line and the keyword
- 0003.3
  - Given Developer has an AND line no deeper than the step it extends
    - When `requirements` reads it
      - Then it fails naming the line
- 0003.4
  - Given Developer has a Gherkin step that holds a comma outside code
    - When `requirements` reads it
      - Then it fails asking for an AND line
- 0003.5
  - Given Developer has a Gherkin block the Cucumber parser rejects
    - When `requirements` reads it
      - Then it fails at the doc line the parser names
- 0003.6
  - Given Developer has a doc of bullet acceptances
    1. When migration `0003.6-gherkin` runs up through `--fix` or `--up`
       - Then each acceptance becomes an all-caps SCENARIO in one Gherkin block that passes
    2. When it runs down through `--down`
       - Then the doc holds bullet acceptances again with each AND joined to its line
- 0003.7
  - Given Developer moves a bullet doc to Gherkin and splits its commas into AND lines
    - When `requirements` judges the change
      - Then the move pairs as a reformat so its tests need no change
