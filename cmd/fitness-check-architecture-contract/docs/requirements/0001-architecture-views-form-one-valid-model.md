# 0001 Architecture Views Form One Valid Model

## Why

I get my five architecture views read as one model, with every broken declaration named at its line.

## Measurement

repos whose architecture views form one valid model

-

repos with architecture views or product code

## Requirements

```gherkin
SCENARIO: 0001.1
  GIVEN Developer has a repo with or without product code AND architecture views
  WHEN it has no views
  THEN product code fails asking for all five
    BUT tests pass with no files
    AND tool configs pass with no files
    AND generated files pass with no files
  WHEN a view is missing
    AND a view has no diagram outside its legends
  THEN it fails naming the view
  WHEN the views sit under `architecture/` or `docs/architecture/` or both
  THEN they pass as one model
    AND it leaves out legends
    AND it leaves out styling lines
    AND it leaves out unclaimed sources
  WHEN callouts are renumbered
  THEN diagram IDs keep their identity
    AND it still passes
  WHEN a pull request changes only tests or changelog or docs
  THEN it passes with no files
  WHEN a pull request changes another file while a view stays broken
  THEN the broken view stays reported

SCENARIO: 0001.2
  GIVEN Developer has five views with diagrams AND callout tables
  WHEN IDs repeat
    AND IDs are invalid
    AND edges name undeclared IDs
  THEN it fails naming the view line
    AND it names the first line of a duplicate
  WHEN a container or component sits outside any boundary
  THEN it fails naming its line
  WHEN one person ID has two names
    AND one name has two IDs
  THEN it fails naming both lines
  WHEN a claimed path or `*` pattern is missing
    AND two elements claim one path or tracked root source
  THEN it fails naming the claim lines
    BUT one element may enclose the other
  WHEN an evidence link names a missing file or line range
  THEN it fails naming the row
  WHEN a requirement reference resolves to no doc or acceptance or chained step
  THEN it fails naming the line
    AND it leaves doc shape to the `requirements` check
  WHEN a diagram uses an unsupported Mermaid form or C4 macro
  THEN it fails naming the line
```
