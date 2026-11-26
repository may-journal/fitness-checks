# mermaid-callouts

Requires every mermaid diagram but a legend to carry callouts, and a callout table whose callouts match it one to one. Orphan callouts, orphan rows, and duplicates fail. It is the foundational check of the mermaid diagram and callout table fitness-function set ([#28](https://github.com/may-journal/fitness-runner/issues/28)).

## Behavior

[Requirement 0001](docs/requirements/0001-mermaid-diagrams-carry-matching-callouts.md) states what this check does, and a test case named for each step proves it.

## Reading callouts

A node callout is the leading integer of a node label: quoted (`Container(x, "3", …)`, `A["1 Check"]`) or inside a node shape (`A[2 Screen]`). A wardley-beta map numbers its annotations instead, as in `annotation 1,[0.40, 0.35] "text"`.

A connector callout is the leading label of an edge, in any of these forms:

- C4 relationship: `Rel(a, b, "A")`.
- Flowchart pipe: `a -->|B| b` or `a -->|"B"| b`.
- Inline link: `a -- C --> b`.
- Colon after an arrow: `a --> b : D` or `a ->> b: D`.

A colon label with no arrow before it, such as a state description or a note, is a node's. Capital letters that break the doubling count only when they are the whole label, so prose such as `|API calls|` is not a callout. Such a label also never makes a diagram need a table on its own. In the table, callouts come from the first (`#`) column: a number, capital letters, or either in brackets.
