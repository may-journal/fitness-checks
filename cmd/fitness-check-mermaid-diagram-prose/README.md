# mermaid-diagram-prose

Fails edge labels that carry prose beyond a callout, number or letter, in any mermaid diagram. Diagram labels must not duplicate callout table prose. It is part of the mermaid diagram and callout table fitness-function set ([#28](https://github.com/may-journal/fitness-runner/issues/28)).

## Behavior

In every diagram, with a callout table or not, it flags edge labels with prose after a callout number or letter. For example, `"6 Uses"` should be `"6"` and `"A calls"` should be `"A"`. Arrows carry callouts only; descriptions belong in the table, not the diagram.

A file with neither a mermaid diagram nor a callout table is skipped and doesn't count toward `filesChecked`. Run `fitness explain mermaid-diagram-prose` for the files it reads.
