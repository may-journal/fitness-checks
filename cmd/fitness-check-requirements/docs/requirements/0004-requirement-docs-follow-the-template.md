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
      1. WHEN a file name and its title ID differ
         1. THEN it fails telling the file name to use
      2. WHEN two docs share an ID
         1. THEN it fails naming the docs
      3. WHEN an acceptance ID repeats in its doc
         1. THEN it fails naming the ID
      4. WHEN a doc holds more acceptances than the cap allows
         1. THEN it fails naming the count
         2. THEN a doc at the limit passes
2. Docs hold the three sections
   1. GIVEN Developer has a repo with requirement docs
      1. WHEN a section is extra or missing or empty
         1. THEN it fails naming the sections
      2. WHEN the Measurement is not a numerator and a dash and a denominator
         1. THEN it fails naming the three lines
      3. WHEN a doc keeps a `Source:` line
         1. THEN it fails naming the line
         2. THEN `fitness requirements --fix` deletes it
3. Acceptances read in order
   1. GIVEN Developer has an acceptance with an extra item or items out of order
      1. WHEN `requirements` reads it
         1. THEN it fails naming the order
4. Step chains stay numbered
   1. GIVEN Developer has acceptances that chain WHENs as steps
      1. WHEN a chain holds up to the step limit beside a single step
         1. THEN it passes
         2. THEN a longer chain fails naming the limit
      2. WHEN a step is misnumbered or unnumbered or the only one numbered
         1. THEN it fails naming the line
      3. WHEN two docs share a WHEN
         1. THEN it fails naming both docs
5. IDs are never reused
   1. GIVEN Developer has a repo whose history removed an ID
      1. WHEN a commit brings a removed ID back
         1. THEN it fails naming the ID
      2. WHEN a commit brings a removed step back
         1. THEN it fails naming the step
      3. WHEN a commit removes every doc
         1. THEN the history restarts
         2. THEN later removals count afresh
