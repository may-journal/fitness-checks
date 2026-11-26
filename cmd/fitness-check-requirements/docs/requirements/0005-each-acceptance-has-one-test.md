# 0005 Each Acceptance Has One Test

## Why

I see every acceptance proven by exactly one test named for it, and every test proving a requirement.

## Measurement

acceptances proven by exactly one test

-

acceptances

## Requirements

1. Each acceptance has its one test
   1. GIVEN Developer has a repo with requirement docs plus tests
      1. WHEN an acceptance has no test
         1. THEN it fails naming the test to write
      2. WHEN two tests own one acceptance
         1. THEN it fails naming the count
      3. WHEN a test is named for an acceptance no doc defines
         1. THEN it fails naming the test
         2. THEN it fails naming the ID
2. Each test proves a requirement
   1. GIVEN Developer has a test named for no acceptance or for a step
      1. WHEN `requirements` reads the tests
         1. THEN a test proving no requirement fails naming it
      2. WHEN a test is named for a step
         1. THEN it fails telling to name the acceptance's test
3. Each chained step has a case
   1. GIVEN Developer has an acceptance that chains steps
      1. WHEN each step has a case named with its ID
         1. THEN it passes beside cases named otherwise
      2. WHEN a step has no case
         1. THEN it fails naming the step
         2. THEN it fails naming the test to add it to
      3. WHEN a case names a step its acceptance does not chain
         1. THEN it fails naming the case
      4. WHEN the chain is malformed
         1. THEN its good cases are not called stray
