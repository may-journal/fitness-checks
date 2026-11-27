# 0012 Describe Answers Are Compared

## Why

I see whether a change altered how the engine calls a check, so a batch of docs and tests moves no check's answer.

## Measurement

checks whose `--describe` answer matches their last release

-

checks compared

## Requirements

1. Describe answers are compared
   1. GIVEN Developer has a check released at a tag before a commit
      1. WHEN its `--describe` answer is unchanged since the tag
         1. THEN the comparison prints identical
         2. THEN it exits 0
      2. WHEN its `--describe` answer changed since the tag
         1. THEN the comparison prints the diff
         2. THEN it exits 1
