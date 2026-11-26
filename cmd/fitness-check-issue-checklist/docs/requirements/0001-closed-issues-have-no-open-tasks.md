# 0001 Closed Issues Have No Open Tasks

## Why

I never close an issue whose task list still holds unfinished work.

## Measurement

issue bodies whose task items are all ticked

-

issue bodies with a task list

## Requirements

1. Unchecked task items fail with their text
   1. GIVEN Developer has an issue body
   2. GIVEN it may hold a task list
      1. WHEN the body has unchecked items
         1. THEN it fails with one error per item quoting its text
      2. WHEN every item is ticked
         1. THEN it passes
         2. THEN a body with no task list passes
      3. WHEN an unchecked item sits inside fenced code
         1. THEN it never counts since it is an example
      4. WHEN there is no input
         1. THEN it passes with zero files checked
