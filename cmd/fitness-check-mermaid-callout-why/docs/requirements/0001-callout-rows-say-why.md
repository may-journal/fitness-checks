# 0001 Callout Rows Say Why

## Why

I learn why each node and connector exists from its row. A table without a Why column or with an empty Why cell tells me nothing.

## Measurement

Callout tables whose numbered rows each give a reason

-

Callout tables

## Requirements

1. Callout rows give a Why
   1. GIVEN Developer has Markdown files or a PR body with callout tables
      1. WHEN a table has no Why column
         1. THEN it fails naming the table line
      2. WHEN a numbered row has an empty Why cell
         1. THEN it fails naming the row and its line
         2. THEN a row with no Why cell fails the same way
      3. WHEN every numbered row has a Why in any header case
         1. THEN the table passes
      4. WHEN an unnumbered row has an empty Why
         1. THEN the row passes
      5. WHEN a PR body table has no Why column
         1. THEN it fails naming `(description)`
      6. WHEN a file has no diagram and no callout table
         1. THEN it is skipped
         2. THEN it is not counted
