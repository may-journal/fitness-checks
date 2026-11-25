# 0001 Step Chains Stay Within The Configured Limit

## Why

My requirement docs chain steps up to my repo's limit, in the repo root and in each package, and other lists stay short.

## Measurement

requirement docs whose step chains pass up to the repo's limit

-

requirement docs with step chains

## Requirements

```gherkin
SCENARIO: 0001.1
  GIVEN Developer sets `requirements.maxSteps` and has a root or package requirement doc that chains steps
  WHEN `prose-budget` reads it
  THEN a chain up to that limit passes and one step more fails

SCENARIO: 0001.2
  GIVEN Developer sets no `requirements.maxSteps` or one below two
  WHEN `prose-budget` reads a step chain in a root or package doc
  THEN the default limit of seven steps applies
```
