# 0004 Requirement Docs Follow The Template

## Why

Every requirement doc has one ID and three sections and ordered acceptances, so I read each repo's docs the same way.

## Measurement

requirement docs that follow the template

-

requirement docs

## Requirements

1. Docs are named and numbered
   1. GIVEN Developer has a repo with requirement docs
      1. WHEN a file name differs from its title ID
         1. THEN it fails telling the file name to use
      2. WHEN two docs share an ID
         1. THEN it fails naming the docs
      3. WHEN two acceptances in a doc share a name
         1. THEN it fails naming both IDs
      4. WHEN a doc holds more acceptances than the cap allows
         1. THEN it fails naming the count
         2. THEN a doc at the limit passes
2. Docs hold the three sections
   1. GIVEN Developer has a repo with requirement docs
      1. WHEN a section is extra or missing or empty
         1. THEN it fails naming the sections
      2. WHEN the Measurement is not numerator over dash over denominator
         1. THEN it fails naming the three lines
      3. WHEN a doc keeps a `Source:` line
         1. THEN it fails naming the line
         2. THEN `fitness requirements --fix` deletes it
3. Removed in 2026.11.26.1608: 0003.1 states order
4. Step chains stay numbered
   1. GIVEN Developer has acceptances that chain WHENs as steps
      1. WHEN a chain holds up to the step limit beside a single step
         1. THEN it passes
         2. THEN a longer chain fails naming the limit
      2. WHEN a step is misnumbered or unnumbered
         1. THEN it fails naming the line
      3. WHEN two docs share a WHEN
         1. THEN it fails naming both docs
5. IDs are never reused
   1. GIVEN Developer has a repo whose history removed an ID
      1. WHEN a removed ID returns in a commit
         1. THEN it fails naming the ID
      2. WHEN a removed step returns in a commit
         1. THEN it fails naming the step
      3. WHEN a commit removes every doc
         1. THEN the history restarts
         2. THEN later removals count afresh
6. Bullets fail
   1. GIVEN Developer has a bullet doc
      1. WHEN `requirements` reads it
         1. THEN it fails at each acceptance
         2. THEN it links the numbered form
         3. THEN its tests keep their IDs
