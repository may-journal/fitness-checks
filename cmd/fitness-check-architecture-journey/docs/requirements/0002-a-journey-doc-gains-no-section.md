# 0002 A Journey Doc Gains No Section

## Why

I extend an existing journey flow for a new requirement, since a new section splits the journey in two.

## Measurement

journey docs whose sections all exist on the base branch

-

journey docs changed in a pull request

## Requirements

1. A journey doc gains no new section
   1. GIVEN Developer has a journey doc on the base branch
   2. GIVEN a pull request into that branch changes it
      1. WHEN a `##` section is added
         1. THEN it fails at the new heading's line
         2. THEN a section with its own diagram fails telling to extend an existing flow
         3. THEN a commit adding one fails at the same line
      2. WHEN no `##` section is added
         1. THEN a node added to an existing diagram passes
         2. THEN a `###` heading passes
         3. THEN a `##` line inside a code fence passes
         4. THEN a journey doc new in the pull request passes
