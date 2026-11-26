# 0001 Components Own Every Source

## Why

I get every product source owned by one drawn component, and every reference between them drawn.

## Measurement

repos whose components view owns and draws their product sources

-

repos with a components view

## Requirements

```gherkin
SCENARIO: 0001.1
  GIVEN Developer has a components view AND product sources
  WHEN a product source is under no component's claimed path
    AND the source may be agent tooling
  THEN it fails naming the source
  WHEN a source is a root tool config
    AND a source is under `scripts/`
    AND a source is run only by manifest scripts
  THEN it needs no owner
    BUT the product may start or publish it
  WHEN a code span or label names a root-level source by bare name
  THEN that component owns it
    BUT other bare names stay prose
  WHEN two diagrams of one view claim the same source
  THEN it fails naming both claims
  WHEN a reference between owners has no edge drawn between them or enclosing boundaries
  THEN it fails naming the view
    AND it names any reverse edge
    BUT unused edges pass
  WHEN a reference sits in a comment
    AND a reference sits in a string
    AND a reference sits in generated source
    AND a reference sits under a build condition
  THEN only the build condition counts
    AND it is named after `under`
  WHEN a product source cannot be parsed
  THEN it fails naming the source

SCENARIO: 0001.2
  GIVEN Developer has a repo with or without a components view
  WHEN the repo has no architecture views
  THEN it passes with zero files
  WHEN a name resolves to a declaration in another component's file
  THEN it is a reference the view must draw
  WHEN an edge row links a source
  THEN the first component must own it
    AND its lines name the second outside comments
  WHEN a runtime reaches a component under another runtime's container boundary
  THEN it fails naming that boundary
    BUT components outside every container are shared
  WHEN a scoped run changes only tests or changelog or docs
  THEN it passes with zero files
  WHEN a scoped run changes a view or a product source
  THEN it reads the whole model
    AND it reads every product source
```
