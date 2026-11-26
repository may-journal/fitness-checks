# 0005 Diagrams Match The Code

## Why

I can trust every diagram I keep. Each one renders, explains its callouts, and adds detail at each level.

## Measurement

diagrams that render and match their callout tables

-

diagrams checked

## Requirements

1. Removed in 2026.11.26.0300: moved to the mermaid-syntax doc
2. Removed in 2026.11.26.0300: moved to the mermaid-callouts doc
3. Removed in 2026.11.26.0300: moved to the mermaid-callouts doc
4. Removed in 2026.11.26.0300: moved to the mermaid-callout-why and mermaid-legend and mermaid-diagram-table-gap and mermaid-diagram-prose docs
5. Removed in 2026.11.25.0800: moved to the `architecture-contract` doc
6. Removed in 2026.11.25.0800: moved to the `architecture-journey` and `architecture-actors` docs
7. Removed in 2026.11.25.0800: moved to the `architecture-containers` and `architecture-components` docs
8. Removed in 2026.11.25.0800: moved to the `architecture-code` doc
9. Levels do not repeat
   1. GIVEN Developer has architecture level docs
   2. GIVEN their callout descriptions repeat a parent level's text
      1. WHEN the Mermaid checks compare each level with the one above it
         1. THEN a verbatim repeat fails naming the file and level
         2. THEN distinct or blank or distant levels pass
