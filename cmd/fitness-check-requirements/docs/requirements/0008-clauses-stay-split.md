# 0008 Clauses Stay Split

## Why

I read one condition or one result in each step, so no step hides a second clause behind a joining word.

## Measurement

steps that state one condition or one result

-

steps in numbered acceptances

## Requirements

1. Clauses stay split
   1. GIVEN Developer has numbered acceptances `requirements` reads
      1. WHEN a step holds a lowercase "and" or "but" outside code or quotes
         1. THEN it asks to split at the word
      2. WHEN the word sits in code or in quotes
         1. THEN it passes
