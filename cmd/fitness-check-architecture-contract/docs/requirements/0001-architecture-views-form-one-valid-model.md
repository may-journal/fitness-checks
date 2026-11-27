# 0001 Architecture Views Form One Valid Model

## Why

I get my five architecture views read as one model, with every broken declaration named at its line.

## Measurement

repos whose architecture views form one valid model

-

repos with architecture views or product code

## Requirements

1. Missing views fail product code
   1. GIVEN Developer has a repo
   2. GIVEN it may hold product code
   3. GIVEN it may hold architecture views
      1. WHEN it has no views
         1. THEN product code fails asking for all five
         2. THEN tests pass with no files
         3. THEN tool configs pass with no files
         4. THEN generated files pass with no files
      2. WHEN a view is missing
         1. THEN it fails naming the view
         2. THEN a view with no diagram outside its legends fails the same way
      3. WHEN some views sit under `architecture/` or `docs/architecture/`
         1. THEN they pass as one model
      4. WHEN callouts are renumbered
         1. THEN diagram IDs keep their identity
      5. WHEN only tests or changelog or docs change in a pull request
         1. THEN it passes with no files
      6. WHEN another file changes in a pull request while a view stays broken
         1. THEN the broken view stays reported
2. The five views agree
   1. GIVEN Developer has five views with diagrams
   2. GIVEN each has callout tables
      1. WHEN IDs repeat or are invalid
         1. THEN it fails naming the view line
         2. THEN an edge naming an undeclared ID fails the same way
      2. WHEN a container or component sits outside any boundary
         1. THEN it fails naming its line
      3. WHEN one person ID has two names
         1. THEN it fails naming both lines
      4. WHEN a claimed path or `*` pattern is missing
         1. THEN it fails naming the claim lines
         2. THEN two elements claiming one path or tracked root source fail the same way
      5. WHEN an evidence link names a missing file or line range
         1. THEN it fails naming the row
      6. WHEN a requirement reference resolves to no doc
         1. THEN it fails naming the reference line
         2. THEN an acceptance or step resolving to nothing fails the same way
      7. WHEN a diagram uses an unsupported Mermaid form or C4 macro
         1. THEN it fails naming the line
