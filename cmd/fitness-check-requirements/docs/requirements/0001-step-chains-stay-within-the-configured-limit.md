# 0001 Step Chains Stay Within The Configured Limit

## Why

I chain a requirement's steps as far as my repo allows and no further, so each acceptance stays one short story.

## Measurement

step chains within the repo's limit

-

step chains

## Requirements

```gherkin
SCENARIO: 0001.1
  GIVEN Developer sets `requirements.maxSteps` and has a doc that chains one step more
  WHEN `requirements` reads it
  THEN it fails naming the chain and the limit

SCENARIO: 0001.2
  GIVEN Developer sets no `requirements.maxSteps` or one below two
  WHEN `requirements` reads a step chain
  THEN the default limit of seven steps applies
```
