# mermaid-callout-why

Requires each numbered callout table to carry a "Why" column, and flags numbered rows whose "Why" cell is empty. It is part of the mermaid diagram and callout table fitness-function set ([#28](https://github.com/may-journal/fitness-runner/issues/28)).

## Behavior

For each numbered callout table, it requires a `Why` column (case-insensitive) in the header. It warns when a callout row leaves its `Why` cell empty. A callout row opens with a node number, such as `1`, or a connector letter, such as `A`, as [mermaid-callouts](../fitness-check-mermaid-callouts/README.md) reads them.

A file with neither a mermaid diagram nor a callout table is skipped and doesn't count toward `filesChecked`. Run `fitness explain mermaid-callout-why` for the files it reads.
