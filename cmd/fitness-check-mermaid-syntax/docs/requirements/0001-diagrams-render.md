# 0001 Diagrams Render

## Why

I never meet a diagram that fails to draw. Each broken one tells me its line and kind.

## Measurement

Mermaid diagrams that parse

-

Mermaid diagrams

## Requirements

1. Each diagram type parses
   1. GIVEN Developer has Markdown files with mermaid diagrams
      1. WHEN every diagram is well formed in a known type
         1. THEN the flowchart passes
         2. THEN the sequence diagram passes
         3. THEN the C4 diagram passes
      2. WHEN a flowchart has an edge with no target or a subgraph with no end
         1. THEN it fails naming each line
         2. THEN it names the diagram line
      3. WHEN a sequence diagram has a message with no target or a block with no end
         1. THEN it fails naming each line
      4. WHEN a diagram type is unknown
         1. THEN it fails naming the type
      5. WHEN a C4 diagram has an unclosed bracket
         1. THEN it fails naming the line
      6. WHEN an unquoted flowchart label holds brackets or quotes
         1. THEN it fails asking for double quotes
2. Diagrams are found wherever they sit
   1. GIVEN Developer has diagrams in many kinds of file
   2. GIVEN some sit in a PR body
      1. WHEN a PR body diagram is broken
         1. THEN it fails naming `(description)`
      2. WHEN no file holds a diagram
         1. THEN the check passes with no files
      3. WHEN a diagram sits in a `.mmd` file or a blockquote or a `.markdown` file
         1. THEN each real diagram is judged
         2. THEN a code block example is not judged
