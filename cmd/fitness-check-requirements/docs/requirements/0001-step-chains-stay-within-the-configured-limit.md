# 0001 Step Chains Stay Within The Configured Limit

## Why

I chain a requirement's steps as far as my repo allows and no further, so each acceptance stays one short story.

## Measurement

step chains within the repo's limit

-

step chains

## Requirements

1. A set limit holds
   1. GIVEN Developer sets `requirements.maxSteps`
   2. GIVEN a doc chains one step more
      1. WHEN `requirements` reads it
         1. THEN it fails naming the chain
         2. THEN it names the limit
2. The default limit is seven
   1. GIVEN Developer sets no `requirements.maxSteps`
   2. GIVEN Developer sets one below two
      1. WHEN `requirements` reads a step chain
         1. THEN the default limit of seven steps applies
