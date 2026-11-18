# 0079 Checks Mermaid Syntax

## Why

Every mermaid diagram I publish renders, because a broken one fails before anyone reads it.

## Measurement

mermaid diagrams that parse
-
mermaid diagrams

## Requirements

- 0079.1
    - Given Developer has valid flowcharts, graphs, sequence diagrams, and other types, with inline edge labels, description blocks, and shapes that leave brackets unpaired
        - When mermaid-syntax checks the doc
            - Then it passes
- 0079.3
    - Given Developer has a flowchart with an edge missing its target and an unclosed subgraph
        - When mermaid-syntax checks the doc
            - Then it fails naming each line
- 0079.4
    - Given Developer has a sequence diagram with a message missing its colon and an unclosed loop
        - When mermaid-syntax checks the doc
            - Then it fails naming each line
- 0079.5
    - Given Developer has a diagram of an unknown type
        - When mermaid-syntax checks the doc
            - Then it fails naming the type
- 0079.6
    - Given Developer has a C4 diagram with an unclosed bracket
        - When mermaid-syntax checks the doc
            - Then it fails naming the bracket's line
- 0079.7
    - Given Developer has a PR body with a broken flowchart
        - When mermaid-syntax checks it in body mode
            - Then it fails like a file would
- 0079.8
    - Given Developer has docs with no mermaid diagrams
        - When mermaid-syntax runs
            - Then it passes having checked no files
- 0079.9
    - Given Developer has bare brackets in a doc, a quote, a quote cut short, a `.markdown` doc, a `.mmd` file and quoted code
        - When mermaid-syntax checks the repo
            - Then it fails naming each real node and its line
