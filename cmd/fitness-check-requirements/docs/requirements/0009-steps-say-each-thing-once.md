# 0009 Steps Say Each Thing Once

## Why

I read each condition once, in the GIVEN, and each result once, so a repeat tells me to group the steps.

## Measurement

numbered acceptances whose WHENs and THENs each say one thing

-

numbered acceptances

## Requirements

1. WHENs that open alike
   1. GIVEN Developer has numbered acceptances `requirements` reads
      1. WHEN sibling WHENs open with the same three meaningful words
         1. THEN it names the shared words to move into the GIVEN
      2. WHEN the words they share are filler or fewer than three
         1. THEN it passes
2. THENs that repeat
   1. GIVEN Developer has numbered acceptances `requirements` reads
      1. WHEN two THENs of one acceptance match apart from case or trailing punctuation
         1. THEN it names both lines to group under one WHEN
         2. THEN it says to drop a THEN repeated under the same WHEN
      2. WHEN the same THEN sits in two acceptances
         1. THEN it passes
