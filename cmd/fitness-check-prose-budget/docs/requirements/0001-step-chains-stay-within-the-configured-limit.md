# 0001 Step Chains Stay Within The Configured Limit

## Why

My requirement docs chain steps up to my repo's limit, in the repo root and in each package, and other lists stay short.

## Measurement

requirement docs whose step chains pass up to the repo's limit

-

requirement docs with step chains

## Requirements

- 0001.1
  - Given Developer sets `requirements.maxSteps` and has a root or package requirement doc that chains steps
    - When `prose-budget` reads it
      - Then a chain up to that limit passes and one step more fails
- 0001.2
  - Given Developer sets no `requirements.maxSteps` or one below two
    - When `prose-budget` reads a step chain in a root or package doc
      - Then the default limit of seven steps applies
- 0001.3
  - Given Developer has a requirement doc of numbered acceptances
    1. When its `Removed` lines pass the list limit
       - Then they count toward no limit
    2. When its GIVENs or a WHEN's THENs pass half the list limit
       - Then it fails at that level
