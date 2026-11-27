# 0001 Code View Matches The Source

## Why

I keep my code view's classes and edges true to the declarations and relations my source holds.

## Measurement

code view classes and edges that match the source

-

code view classes and edges

## Requirements

1. The code view names real symbols
   1. GIVEN Developer has a code view
   2. GIVEN its concrete classes link their declaring sources
      1. WHEN each class matches its source declaration
         1. THEN it passes with generics
         2. THEN it passes with extensions
         3. THEN it passes with any spacing or methods elsewhere
      2. WHEN a class's linked lines lack its symbol
         1. THEN it fails naming its declaration if any
      3. WHEN a kind or signature or member differs from the source
         1. THEN it fails naming the class line
         2. THEN it fails naming the source line
      4. WHEN one name is declared twice
         1. THEN it fails listing both until a link or signature picks one
      5. WHEN a class links a test or generated file
         1. THEN it fails
      6. WHEN two classes name one declaration
         1. THEN it fails unless the class is `<<abstract>>`
         2. THEN a label naming no symbol fails too
      7. WHEN a shown folder holds an exported type no class names
         1. THEN it fails naming the type
2. Edges allow only what the source uses
   1. GIVEN Developer has a code view drawing edges between its classes
      1. WHEN the source uses a class no edge allows
         1. THEN it fails at the source line naming the smallest holding class
      2. WHEN an edge is drawn with an arrow or a line
         1. THEN an arrow allows uses from the first class
         2. THEN a plain line allows uses either way
      3. WHEN a drawn association or realization or inheritance has no source relation
         1. THEN it fails at the edge
      4. WHEN an edge row links a source
         1. THEN the linked lines must hold that relation
      5. WHEN an edge touches an abstract class
         1. THEN it needs no use
      6. WHEN an edge uses any other arrow
         1. THEN it fails as unsupported
      7. WHEN the repo has no views
         1. THEN it passes with no files
         2. THEN a test-only pull request passes too
