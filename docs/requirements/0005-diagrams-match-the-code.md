# 0005 Diagrams Match The Code

## Why

I can trust every diagram I keep. Each one renders, explains its callouts, and matches the code and requirements it draws.

## Measurement

diagrams that render and match their callout tables and the code they describe

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
  - Given Developer has architecture levels and a contract with diagram IDs, source paths and requirement links
    - When the architecture checks read the contract and its levels
      - Then repeated level text, duplicate claims and unreadable evidence fail naming their lines
- 0005.6
  - Given Developer has a journey diagram and a context diagram beside requirement docs and their Givens
    - When the architecture checks compare requirements with journey steps and actors
      - Then each unnamed doc, unfulfilled step, crossed label and undeclared actor fails naming its line
- 0005.7
  - Given Developer has container and component views beside build targets and product folders
    - When the architecture checks compare the views with the runtimes and folders
      - Then each undeclared, stale or unowned container, folder or call fails naming its source
- 0005.8
  - Given Developer has a code view of classes, symbols and edges
    - When the architecture checks compare the code view with the source
      - Then each omitted, stale, duplicate or missing class and edge fails naming its line
