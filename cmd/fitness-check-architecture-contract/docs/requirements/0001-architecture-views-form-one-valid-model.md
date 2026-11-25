# 0001 Architecture Views Form One Valid Model

## Why

I get my five architecture views read as one model, with every broken declaration named at its line.

## Measurement

repos whose architecture views form one valid model

-

repos with architecture views or product code

## Requirements

- 0001.1
  - Given Developer has a repo, with or without product code and architecture views
    1. When it has no views
       - Then product code fails asking for all five, while tests, tool configs and generated files pass with no files
    2. When a view is missing, or one has no diagram outside its legends
       - Then it fails naming the view
    3. When the views sit under `architecture/`, `docs/architecture/` or both
       - Then they pass as one model, leaving out legends, styling lines and unclaimed sources
    4. When callouts are renumbered
       - Then diagram IDs keep their identity and it still passes
    5. When a pull request changes only tests, changelog or docs
       - Then it passes with no files
    6. When a pull request changes another file while a view stays broken
       - Then the broken view stays reported
- 0001.2
  - Given Developer has five views with diagrams and callout tables
    1. When IDs repeat or are invalid, or edges name undeclared IDs
       - Then it fails naming the view line, and the first line of a duplicate
    2. When a container or component sits outside any boundary
       - Then it fails naming its line
    3. When one person ID has two names, or one name two IDs
       - Then it fails naming both lines
    4. When a claimed path or `*` pattern is missing, or two elements claim one path or tracked root source
       - Then it fails naming the claim lines, unless one element encloses the other
    5. When an evidence link names a missing file or line range
       - Then it fails naming the row
    6. When a requirement reference resolves to no doc, acceptance or chained step
       - Then it fails naming the line, leaving doc shape to the `requirements` check
    7. When a diagram uses an unsupported Mermaid form or C4 macro
       - Then it fails naming the line
