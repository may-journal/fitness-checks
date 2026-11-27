# 0010 Requirement Docs Say Each Thing Once

## Why

I read each requirement in this repo once, so no doc repeats a condition or a result.

## Measurement

requirement docs whose acceptances repeat no WHEN opening and no THEN

-

requirement docs

## Requirements

1. Docs repeat nothing
   1. GIVEN Developer has every requirement doc of this repo
      1. WHEN the requirements parser reads each doc
         1. THEN no acceptance holds sibling WHENs that open alike
         2. THEN no acceptance holds a THEN that repeats another
