# 0005 Diagrams Match The Code

## Why

I can trust every diagram I keep. Each one renders, explains its callouts, and matches the code and requirements it draws.

## Measurement

diagrams that render and match their callout tables and the code they describe

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
5. Levels do not repeat
   1. GIVEN Developer has architecture levels
   2. GIVEN a contract holds diagram IDs and source paths
      1. WHEN the architecture checks read the contract
         1. THEN repeated level text fails naming its line
         2. THEN duplicate claims and bad evidence fail
6. Journeys match requirements
   1. GIVEN Developer has journey diagrams beside requirement docs
      1. WHEN the architecture checks compare requirements with steps
         1. THEN an unnamed doc or unfulfilled step fails naming its line
         2. THEN a crossed label fails naming its line
         3. THEN an undeclared actor or stale region fails naming its line
7. Containers match the build
   1. GIVEN Developer has container views beside the build
      1. WHEN the architecture checks compare them
         1. THEN an undeclared or stale container fails naming its source
         2. THEN an unowned folder or undeclared call fails naming its source
8. Code views match source
   1. GIVEN Developer has a code view of classes and symbols and edges
      1. WHEN the architecture checks compare it with the source
         1. THEN an omitted or stale class fails naming its line
         2. THEN a missing edge fails naming its line
