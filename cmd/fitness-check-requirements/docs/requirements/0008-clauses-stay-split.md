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
      1. WHEN a step holds a word of `given when then and but` not in caps outside code or quotes
         1. THEN it asks for caps or a split into its own step
      2. WHEN the word sits in code or in quotes or in caps
         1. THEN it passes
2. Block keywords stay in caps
   1. GIVEN Developer has numbered acceptances `requirements` reads
      1. WHEN a block keyword with its colon opens a line in mixed case
         1. THEN it asks for the keyword in caps
      2. WHEN the keyword is in caps or no colon follows
         1. THEN it passes
