# 0001 Mermaid Diagrams Carry Matching Callouts

## Why

I read each diagram by its numbered nodes and lettered connectors. A callout table explains each one and nothing is left unpaired.

## Measurement

Markdown files whose diagrams and callout tables pair one to one

-

Markdown files with a diagram or a callout table

## Requirements

1. Callouts and rows pair one to one
   1. GIVEN Developer has Markdown files or a PR body with mermaid diagrams
   2. GIVEN each diagram has a callout table
      1. WHEN every callout in a diagram maps to exactly one table row in both directions
         1. THEN the file passes
         2. THEN the count includes the file
      2. WHEN a diagram other than a legend has no numbered callouts
         1. THEN it fails with `mermaid diagram has no numbered callouts`
      3. WHEN a callout does not pair with its row
         1. THEN a callout with no row fails naming it with its line
         2. THEN a row with no callout fails naming it with its line
         3. THEN a callout or row that appears twice fails the same way
         4. THEN a callout with no table or a row with no diagram fails the same way
      4. WHEN a connector callout is a number
         1. THEN it fails asking for a letter
         2. THEN the failure names `--fix`
      5. WHEN a connector letter breaks the doubling from `A` to `Z` or from `AA` to `ZZ`
         1. THEN it fails for the diagram callout
         2. THEN it fails for the table row
      6. WHEN a PR body diagram has a callout with no row
         1. THEN it fails naming `(description)`
2. Fix lettered connectors once
   1. GIVEN Developer has Markdown files with numbered connectors
   2. GIVEN some files hold no diagram
      1. WHEN it runs with `--fix` on numbered connectors
         1. THEN each becomes the next free letter
         2. THEN its row follows
         3. THEN the repo passes
      2. WHEN `--fix` runs a second time
         1. THEN it changes nothing
      3. WHEN a numbered connector shares its number with a node
         1. THEN `--fix` leaves the diagram unchanged
         2. THEN the check still fails
      4. WHEN a file has neither a diagram nor a callout table
         1. THEN it is skipped
         2. THEN it is not counted
