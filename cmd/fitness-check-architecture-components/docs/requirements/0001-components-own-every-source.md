# 0001 Components Own Every Source

## Why

I get every product source owned by one drawn component, and every reference between them drawn.

## Measurement

repos whose components view owns and draws their product sources

-

repos with a components view

## Requirements

- 0001.1
  - Given Developer has a components view and product sources
    1. When a product source, agent tooling included, is under no component's claimed path
       - Then it fails naming the source
    2. When a source is a root tool config, under `scripts/`, or run only by manifest scripts
       - Then it needs no owner, unless the product starts or publishes it
    3. When a code span or label names a root-level source by bare name
       - Then that component owns it; other bare names stay prose
    4. When two diagrams of one view claim the same source
       - Then it fails naming both claims
    5. When a reference between owners has no edge drawn between them or enclosing boundaries
       - Then it fails naming the view and any reverse edge; unused edges pass
    6. When a reference sits in a comment, string, generated source or under a build condition
       - Then only the build condition counts, named after `under`
    7. When a product source cannot be parsed
       - Then it fails naming the source
- 0001.2
  - Given Developer has a repo, with or without a components view
    1. When the repo has no architecture views
       - Then it passes with zero files
    2. When a name resolves to a declaration in another component's file
       - Then it is a reference the view must draw
    3. When an edge row links a source
       - Then the first component must own it and its lines name the second outside comments
    4. When a runtime reaches a component under another runtime's container boundary
       - Then it fails naming that boundary; components outside every container are shared
    5. When a scoped run changes only tests, changelog or docs
       - Then it passes with zero files
    6. When a scoped run changes a view or a product source
       - Then it reads the whole model and every product source
