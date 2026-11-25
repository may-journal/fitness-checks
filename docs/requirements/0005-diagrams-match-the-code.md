# 0005 Diagrams Match The Code

## Why

I can trust every diagram I keep. Each one renders, explains its callouts, and adds detail at each level.

## Measurement

diagrams that render and match their callout tables

-

diagrams checked

## Requirements

- 0005.1
  - Given Developer has Mermaid diagrams of any type in docs, PR bodies and `.mmd` files
    - When the Mermaid checks parse each diagram
      - Then every broken diagram fails naming its line and type
- 0005.2
  - Given Developer has diagrams whose callouts and callout table rows should pair one to one
    - When the Mermaid checks match callouts to table rows
      - Then each missing, repeated or unpaired callout fails naming its line
- 0005.3
  - Given Developer has connectors labelled with numbers or out-of-order letters
    - When the Mermaid checks read connector labels and run with `--fix`
      - Then bad letters fail, and a fix gives each connector and row the next free letter once
- 0005.4
  - Given Developer has diagrams with callout tables, legends, captions and edge labels
    - When the Mermaid checks read the diagram and its table
      - Then each missing Why, missing legend, extra paragraph or prose label fails naming its line
- 0005.5
  - Given Developer has architecture level docs whose callout descriptions repeat a parent level's text
    - When the Mermaid checks compare each level with the one above it
      - Then a verbatim repeat fails naming the file and level, and distinct, blank or distant levels pass
