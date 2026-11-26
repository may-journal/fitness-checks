# 0001 Components Own Every Source

## Why

I get every product source owned by one drawn component, and every reference between them drawn.

## Measurement

repos whose components view owns and draws their product sources

-

repos with a components view

## Requirements

1. Components own every source
   1. GIVEN Developer has a components view
   2. GIVEN it has product sources
      1. WHEN a product source sits under no component's claimed path
         1. THEN it fails naming the source
         2. THEN agent tooling counts as product source
      2. WHEN a source is a root tool config
         1. THEN it needs no owner
         2. THEN a `scripts/` file needs none
      3. WHEN a code span or label names a root-level source by bare name
         1. THEN that component owns it
      4. WHEN two diagrams of one view claim the same source
         1. THEN it fails naming both claims
      5. WHEN a reference between owners has no edge drawn between them
         1. THEN it fails naming the view
         2. THEN it fails naming any reverse edge
      6. WHEN a reference sits in a comment or string or generated source
         1. THEN it does not count
         2. THEN only a build condition counts
         3. THEN that condition is named after `under`
      7. WHEN a product source cannot be parsed
         1. THEN it fails naming the source
2. Edges follow the references
   1. GIVEN Developer has a repo
   2. GIVEN it may hold a components view
      1. WHEN the repo has no architecture views
         1. THEN it passes with zero files
      2. WHEN a name resolves to a declaration in another component's file
         1. THEN it is a reference the view must draw
      3. WHEN an edge row links a source
         1. THEN the first component must own it
         2. THEN its lines name the second outside comments
      4. WHEN a runtime reaches a component under another runtime's container boundary
         1. THEN it fails naming that boundary
         2. THEN components outside every container are shared
      5. WHEN a scoped run changes only tests or changelog or docs
         1. THEN it passes with zero files
      6. WHEN a scoped run changes a view or a product source
         1. THEN it reads the whole model
         2. THEN it reads every product source
