# 0001 Diagram Labels Hold Callouts Only

## Why

I read a diagram as a map and find every description in the callout table. Edge labels carry only a number or a letter.

## Measurement

Mermaid diagrams whose edge labels are bare callouts

-

Mermaid diagrams

## Requirements

1. Edge labels carry no prose
   1. GIVEN Developer has Markdown files or a PR body with mermaid diagrams
      1. WHEN a C4 relationship label has prose after its callout
         1. THEN it fails naming the label
         2. THEN it names the diagram line
      2. WHEN a flowchart pipe label has prose after its callout
         1. THEN it fails naming each label
      3. WHEN a colon edge label has prose after its callout
         1. THEN it fails naming each label
      4. WHEN every edge label is a bare number or letter
         1. THEN the diagram passes
      5. WHEN a diagram has no callout table
         1. THEN a prose label still fails
2. Prose belongs in node names
   1. GIVEN Developer has Markdown files or a PR body with mermaid diagrams
   2. GIVEN some repos hold no diagram
      1. WHEN prose sits in a node name
         1. THEN the diagram passes
      2. WHEN a PR body diagram has a prose label
         1. THEN it fails naming the label
      3. WHEN a repo has no diagram
         1. THEN the check passes with no files
