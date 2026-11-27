# 0001 Journey And Requirement Docs Cover Each Other

## Why

I get every requirement doc placed on a journey step, and every behavior step naming its doc.

## Measurement

requirement docs a journey diagram label names

-

requirement docs

## Requirements

1. Every doc ID appears on the journey
   1. GIVEN Developer has requirement docs
   2. GIVEN the repo has a journey view or none
      1. WHEN node or subgraph or edge labels name every doc ID on one or more steps
         1. THEN it passes
      2. WHEN a doc ID sits only in tables or prose
         1. THEN it fails at the diagram's fence line
      3. WHEN a label names an acceptance or an unknown doc or step
         1. THEN it fails at its line
         2. THEN a label over eight IDs fails too
      4. WHEN a non-persona node names no ID in its label or enclosing subgraph
         1. THEN it fails at the node's line
      5. WHEN two steps' labels hold each other's doc title words
         1. THEN both fail
      6. WHEN the repo has no views
         1. THEN it passes with no files
         2. THEN a pull request changing only tests passes too
      7. WHEN a pull request changes another file while a doc stays missing
         1. THEN it stays reported
2. Journey nodes follow the docs
   1. GIVEN Developer has a journey view with `journey-nodes` markers listing doc IDs
      1. WHEN each region holds one node per listed doc numbered on from 2
         1. THEN it passes
      2. WHEN a region's lines differ from its listed docs
         1. THEN it fails at its start marker naming `fitness architecture-journey --fix`
      3. WHEN a doc is listed in no region
         1. THEN it fails at the marker
         2. THEN a doc listed in two regions fails the same way
      4. WHEN a start marker lacks its end
         1. THEN it fails at that marker's line
      5. WHEN Developer runs `--fix`
         1. THEN it rewrites every region from the docs
         2. THEN it leaves other lines alone
      6. WHEN a doc lacks its title during `--fix`
         1. THEN it writes nothing
