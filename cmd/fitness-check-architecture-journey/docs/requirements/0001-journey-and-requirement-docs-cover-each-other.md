# 0001 Journey And Requirement Docs Cover Each Other

## Why

I get every requirement doc placed on a journey step, and every behavior step naming its doc.

## Measurement

requirement docs a journey diagram label names

-

requirement docs

## Requirements

- 0001.1
  - Given Developer has requirement docs and a journey view, or neither
    1. When node, subgraph or edge labels name every doc ID, each on one or more steps
       - Then it passes, though a character code like `#9829;` names no doc
    2. When a doc ID sits only in tables or prose
       - Then it fails at the diagram's fence line
    3. When a label names an acceptance, unknown doc or step, or over eight IDs
       - Then it fails at its line, while a cited step reads as its doc
    4. When a non-persona node names no ID in its label or enclosing subgraph
       - Then it fails at its line, since an edge label covers no node
    5. When two steps' labels hold each other's doc title words
       - Then both fail; one alone adds a note and passes
    6. When the repo has no views, or a pull request changes only tests, changelog or docs
       - Then it passes with no files
    7. When a pull request changes another file while a doc stays missing
       - Then it stays reported
- 0001.2
  - Given Developer has a journey view with `journey-nodes` markers listing doc IDs
    1. When each region holds one node per listed doc, numbered on from 2
       - Then it passes, and a view without markers stays hand-written
    2. When a region's lines differ from its listed docs
       - Then it fails at its start marker naming `fitness architecture-journey --fix`
    3. When a doc is listed in no region or two, or a listed ID is no doc
       - Then it fails at the marker
    4. When a start marker lacks its end or nests, or an end has no start
       - Then it fails at that marker's line
    5. When Developer runs `--fix`
       - Then it rewrites every region from the docs, leaving other lines alone
    6. When a doc lacks its title during `--fix`
       - Then it writes nothing and notes the doc
