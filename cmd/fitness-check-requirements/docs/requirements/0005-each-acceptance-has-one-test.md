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
3. Removed in 2026.11.26.2100: each THEN has its own case, which 0005.4 states
4. Each THEN has its own case
   1. GIVEN Developer has an acceptance whose test runs named cases
      1. WHEN each THEN has a case named with its ID then its words
         1. THEN it passes beside cases named otherwise
      2. WHEN a THEN has no case
         1. THEN it fails naming the THEN
         2. THEN it fails naming the test to add the case to
      3. WHEN a case opens with a THEN ID but words it differently
         1. THEN it fails naming the case and the name it needs
      4. WHEN two cases open with one THEN ID
         1. THEN it fails naming the count
      5. WHEN a case names a path no doc holds
         1. THEN it fails naming the case
      6. WHEN a case holds no assertion call
         1. THEN it fails naming the case
      7. WHEN a Node test runs its cases
         1. THEN it needs only its own test until the Node switch is on
