# 0001 No Prose Between Diagram And Table

## Why

I find all detail in the callout table. No paragraph between a diagram and its table repeats it.

## Measurement

Diagrams with no prose before their callout table

-

Diagrams with a callout table

## Requirements

1. The gap holds no paragraph
   1. GIVEN Developer has Markdown files or a PR body with diagrams
   2. GIVEN each diagram has a callout table
      1. WHEN a paragraph sits in the gap from a diagram or a legend to its table
         1. THEN it fails naming the line
      2. WHEN the only line between them is a caption starting with `Numbers` plus the claim they `match the callout table` or prose sits outside the gap
         1. THEN the diagram passes
      3. WHEN a paragraph follows the caption
         1. THEN it fails naming that line
         2. THEN the caption is not named
      4. WHEN a file has no diagram
         1. THEN it is skipped
         2. THEN it is not counted
      5. WHEN a PR body has a paragraph in the gap
         1. THEN it fails naming `(description)`
