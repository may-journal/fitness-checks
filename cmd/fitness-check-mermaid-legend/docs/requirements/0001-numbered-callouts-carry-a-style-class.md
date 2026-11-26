# 0001 Numbered Callouts Carry A Style Class

## Why

I match a numbered node to its legend by its style. Every numbered node carries a class and the diagram declares the classes.

## Measurement

Mermaid diagrams whose numbered nodes carry a style class and a legend

-

Mermaid diagrams with numbered callouts

## Requirements

1. Numbered nodes carry a style
   1. GIVEN Developer has Markdown files or a PR body with flowchart diagrams
      1. WHEN a numbered callout node has no style class
         1. THEN it fails naming the callout
         2. THEN it names the diagram line
      2. WHEN a diagram has numbered callouts and no `classDef` legend
         1. THEN it fails naming the diagram line
      3. WHEN a node is styled by `:::` or `class` or `style`
         1. THEN the diagram passes
      4. WHEN a diagram is not a flowchart
         1. THEN it passes untouched
      5. WHEN a label has no separate number
         1. THEN the label passes
      6. WHEN one of two diagrams has an unstyled callout
         1. THEN it fails naming that diagram line only
2. Diagrams outside files stay judged
   1. GIVEN Developer has diagrams in a PR body
   2. GIVEN some repos hold no diagram
      1. WHEN a PR body callout has no style class
         1. THEN it fails naming `(description)`
      2. WHEN no file holds a diagram
         1. THEN the check passes with no files
