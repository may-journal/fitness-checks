# 0010 Removed Lines Age Out

## Why

I keep docs short. A `Removed` line stops holding its number once two releases pass, and each later test moves up with its acceptance.

## Measurement

Removed lines dropped with their tests renamed

-

Removed lines dropped

## Requirements

1. Old lines drop
   1. GIVEN Developer has a `Removed` line
   2. GIVEN two released versions in the base changelog are newer than its version
      1. WHEN its line is deleted with later acceptances moved up one number
         1. THEN it passes with every moved test case renamed
2. Moved tests are renamed
   1. GIVEN Developer deleted an old `Removed` line
   2. GIVEN a later acceptance moved up one number
      1. WHEN its test keeps the old number
         1. THEN it fails naming the missing test
      2. WHEN its step case keeps the old number
         1. THEN it fails naming the missing case
3. Young lines stay
   1. GIVEN Developer has a `Removed` line
   2. GIVEN fewer than two released versions in the base changelog are newer than its version
      1. WHEN its line is deleted
         1. THEN it fails naming the `Removed` line
         2. THEN its message says how to drop old markers
         3. THEN its message says to rename their tests by hand
