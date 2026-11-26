# 0005 Diagrams Match The Code

## Why

I can trust every diagram I keep. Each one renders, explains its callouts, and adds detail at each level.

## Measurement

diagrams that render and match their callout tables

-

diagrams checked

## Requirements

1. Diagrams parse
   1. GIVEN Developer has Mermaid diagrams in docs and PR bodies
      1. WHEN the Mermaid checks parse each diagram
         1. THEN a broken one fails naming its line
2. Callouts pair with rows
   1. GIVEN Developer has diagrams with callouts and table rows
      1. WHEN the Mermaid checks match callouts to rows
         1. THEN a missing or repeated callout fails naming its line
3. Letters run in order
   1. GIVEN Developer has connectors with numbers or letters out of order
      1. WHEN the Mermaid checks read connector labels
         1. THEN bad letters fail
         2. THEN `--fix` gives each the next free letter
4. Tables are complete
   1. GIVEN Developer has diagrams with callout tables and legends
      1. WHEN the Mermaid checks read the diagram and its table
         1. THEN a missing Why or legend fails
         2. THEN an extra paragraph fails
5. Removed in 2026.11.25.0800: moved to the `architecture-contract` doc
6. Removed in 2026.11.25.0800: moved to the `architecture-journey` and `architecture-actors` docs
7. Removed in 2026.11.25.0800: moved to the `architecture-containers` and `architecture-components` docs
8. Removed in 2026.11.25.0800: moved to the `architecture-code` doc
9. Levels do not repeat
   1. GIVEN Developer has architecture level docs
   2. GIVEN their callout descriptions repeat a parent level's text
      1. WHEN the Mermaid checks compare each level with the one above it
         1. THEN a verbatim repeat fails naming the file and level
         2. THEN distinct or blank or distant levels pass
