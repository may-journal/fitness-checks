# 0001 Journey And Requirement Docs Cover Each Other

## Why

I get every requirement doc placed on a journey step, and every behavior step naming its doc.

## Measurement

requirement docs a journey diagram label names

-

requirement docs

## Requirements

```gherkin
SCENARIO: 0001.1
  GIVEN Developer has requirement docs AND a journey view or neither
  WHEN node or subgraph or edge labels name every doc ID each on one or more steps
  THEN it passes
    BUT a character code like `#9829;` names no doc
  WHEN a doc ID sits only in tables or prose
  THEN it fails at the diagram's fence line
  WHEN a label names an acceptance
    AND a label names an unknown doc
    AND a label names an unknown step
    AND a label names over eight IDs
  THEN it fails at its line
    BUT a cited step reads as its doc
  WHEN a non-persona node names no ID in its label or enclosing subgraph
  THEN it fails at its line
    AND an edge label covers no node
  WHEN two steps' labels hold each other's doc title words
  THEN both fail
    BUT one alone adds a note and passes
  WHEN the repo has no views
    AND a pull request changes only tests or changelog or docs
  THEN it passes with no files
  WHEN a pull request changes another file while a doc stays missing
  THEN it stays reported

SCENARIO: 0001.2
  GIVEN Developer has a journey view with `journey-nodes` markers listing doc IDs
  WHEN each region holds one node per listed doc numbered on from 2
  THEN it passes
    BUT a view without markers stays hand-written
  WHEN a region's lines differ from its listed docs
  THEN it fails at its start marker naming `fitness architecture-journey --fix`
  WHEN a doc is listed in no region
    AND a doc is listed in two regions
    AND a listed ID is no doc
  THEN it fails at the marker
  WHEN a start marker lacks its end
    AND a start marker nests
    AND an end has no start
  THEN it fails at that marker's line
  WHEN Developer runs `--fix`
  THEN it rewrites every region from the docs
    BUT it leaves other lines alone
  WHEN a doc lacks its title during `--fix`
  THEN it writes nothing
    AND it notes the doc
```
