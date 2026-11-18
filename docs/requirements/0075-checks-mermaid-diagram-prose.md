# 0075 Checks Mermaid Diagram Prose

## Why

I read each arrow's meaning once, in the callout table, so a diagram never repeats or contradicts it.

## Measurement

diagrams whose relationship labels hold only callouts
-
diagrams checked

## Requirements

- 0075.1
    - Given Developer has a paired diagram whose `Rel` label has prose after its number or letter
        - When mermaid-diagram-prose checks the doc
            - Then it fails quoting the label
- 0075.2
    - Given Developer has a paired flowchart whose pipe edge label has prose after its number or letter
        - When mermaid-diagram-prose checks the doc
            - Then it fails quoting the label
- 0075.3
    - Given Developer has a paired diagram whose colon edge label has prose after its number or letter
        - When mermaid-diagram-prose checks the doc
            - Then it fails quoting the label
- 0075.4
    - Given Developer has a paired diagram with only bare number labels, comments, and style colors
        - When mermaid-diagram-prose checks the doc
            - Then it passes
- 0075.9
    - Given Developer has a diagram with no callout table whose edge label has prose after its number or letter
        - When mermaid-diagram-prose checks the doc
            - Then it fails quoting the label, as in a paired diagram
- 0075.6
    - Given Developer has a paired diagram whose node names carry prose
        - When mermaid-diagram-prose checks the doc
            - Then it passes, since only arrows are judged
- 0075.7
    - Given Developer has a PR body whose paired diagram has a prose label
        - When mermaid-diagram-prose checks it in body mode
            - Then it fails like a file would
- 0075.8
    - Given Developer has a repo with no diagram or callout table
        - When mermaid-diagram-prose runs
            - Then it passes having judged 0 files
