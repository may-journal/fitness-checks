# 0001 Step Chains Stay Within The Configured Limit

## Why

My requirement docs chain steps up to my repo's limit, in the repo root and in each package, and other lists stay short.

## Measurement

requirement docs whose step chains pass up to the repo's limit

-

requirement docs with step chains

## Requirements

1. A set limit holds
   1. GIVEN Developer sets `requirements.maxSteps`
   2. GIVEN a root or package requirement doc chains steps
      1. WHEN `prose-budget` reads it
         1. THEN a chain up to that limit passes
         2. THEN one step more fails
2. The default limit is seven
   1. GIVEN Developer sets no `requirements.maxSteps`
   2. GIVEN Developer sets one below two
      1. WHEN `prose-budget` reads a step chain in a root or package doc
         1. THEN the default limit of seven steps applies
3. Numbered docs keep lists short
   1. GIVEN Developer has a requirement doc of numbered acceptances
      1. WHEN its `Removed` lines pass the list limit
         1. THEN they count toward no limit
      2. WHEN its GIVENs or a WHEN's THENs pass half the list limit
         1. THEN it fails at that level
4. Requirement docs chain steps by place
   1. GIVEN Developer has a requirement doc that chains numbered steps
      1. WHEN `prose-budget` reads the doc
         1. THEN seven steps pass in a requirement doc
         2. THEN eight steps fail naming the limit
         3. THEN a README keeps the plain list limit
         4. THEN a plain list in a requirement doc keeps half the list limit
