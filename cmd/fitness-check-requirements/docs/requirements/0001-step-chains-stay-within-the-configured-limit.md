# 0001 Step Chains Stay Within The Configured Limit

## Why

I chain a requirement's steps as far as my repo allows and no further, so each acceptance stays one short story.

## Measurement

step chains within the repo's limit

-

step chains

## Requirements

- 0001.1
  - Given Developer sets `requirements.maxSteps` and has a doc that chains one step more
    - When `requirements` reads it
      - Then it fails naming the chain and the limit
- 0001.2
  - Given Developer sets no `requirements.maxSteps` or one below two
    - When `requirements` reads a step chain
      - Then the default limit of seven steps applies
