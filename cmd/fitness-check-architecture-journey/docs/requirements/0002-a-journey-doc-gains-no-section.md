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
      1. WHEN the pull request adds a `##` section
         1. THEN it fails at the new heading's line
      2. WHEN the pull request adds a `##` section with its own diagram
         1. THEN it fails telling to extend an existing flow
      3. WHEN the pull request adds a node to an existing diagram
         1. THEN it passes
      4. WHEN the pull request adds a `###` heading
         1. THEN it passes
      5. WHEN the pull request adds a `##` line inside a code fence
         1. THEN it passes
      6. WHEN the journey doc is new in the pull request
         1. THEN it passes
      7. WHEN Developer commits a new `##` section
         1. THEN it fails at the new heading's line
