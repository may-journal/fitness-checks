# 0005 Diagrams Match The Code

## Why

I can trust every diagram I keep. Each one renders, explains its callouts, and matches the code and requirements it draws.

## Measurement

diagrams that render and match their callout tables and the code they describe

-

diagrams checked

## Requirements

```gherkin
SCENARIO: 0005.1
  GIVEN Developer has Mermaid diagrams of any type in docs and PR bodies and `.mmd` files
  WHEN the Mermaid checks parse each diagram
  THEN every broken diagram fails naming its line and type

SCENARIO: 0005.2
  GIVEN Developer has diagrams whose callouts and callout table rows should pair one to one
  WHEN the Mermaid checks match callouts to table rows
  THEN each missing or repeated or unpaired callout fails naming its line

SCENARIO: 0005.3
  GIVEN Developer has connectors labelled with numbers or out-of-order letters
  WHEN the Mermaid checks read connector labels and run with `--fix`
  THEN bad letters fail
    AND a fix gives each connector and row the next free letter once

SCENARIO: 0005.4
  GIVEN Developer has diagrams with callout tables and legends and captions and edge labels
  WHEN the Mermaid checks read the diagram and its table
  THEN each missing Why or missing legend or extra paragraph or prose label fails naming its line

SCENARIO: 0005.5
  GIVEN Developer has architecture levels and a contract with diagram IDs or source paths and requirement or step links
  WHEN the architecture checks read the contract and its levels
  THEN repeated level text and duplicate claims and unreadable evidence fail naming their lines

SCENARIO: 0005.6
  GIVEN Developer has journey and context diagrams beside bullet or Gherkin docs
  WHEN the architecture checks compare requirements with steps and actors
  THEN each unnamed doc or unfulfilled step or crossed label or undeclared actor or stale region fails naming its line

SCENARIO: 0005.7
  GIVEN Developer has container and component views beside build targets and product folders
  WHEN the architecture checks compare the views with the runtimes and folders
  THEN each undeclared or stale or unowned container or folder or call fails naming its source

SCENARIO: 0005.8
  GIVEN Developer has a code view of classes and symbols and edges
  WHEN the architecture checks compare the code view with the source
  THEN each omitted or stale or duplicate or missing class and edge fails naming its line
```
