# 0001 Code View Matches The Source

## Why

I keep my code view's classes and edges true to the declarations and relations my source holds.

## Measurement

code view classes and edges that match the source

-

code view classes and edges

## Requirements

- 0001.1
  - Given Developer has a code view whose concrete classes link their declaring sources
    1. When each label, link, kind, signature and member matches
       - Then it passes, with generics, extensions, any spacing, and methods elsewhere
    2. When a class's linked lines lack its symbol
       - Then it fails naming where, if anywhere, it is declared
    3. When a kind, interface mark, signature, member or field type differs
       - Then it fails naming the class line and the source line
    4. When one name is declared twice, under build conditions or as overloads
       - Then it fails listing both until a link or signature picks one
    5. When a class links a test, generated or non-source file
       - Then it fails; a non-source link is unsupported evidence
    6. When two classes name one declaration, or a label names no symbol
       - Then it fails, unless the class is `<<abstract>>`
    7. When a shown folder holds an exported type no class names
       - Then it fails; functions and fields may be omitted
- 0001.2
  - Given Developer has a code view drawing edges between its classes
    1. When the source uses a class no edge allows, even via helpers
       - Then it fails at the source line, naming the smallest holding class
    2. When the edge is `-->`, `..>`, `*--`, `o--`, `--` or `..`
       - Then it allows uses from the first class, or either way undirected
    3. When a drawn association, realization or inheritance has no source relation
       - Then it fails at the edge
    4. When an edge row links a source
       - Then the linked lines must hold that relation
    5. When an edge touches an abstract class
       - Then it needs no use; a row link is unsupported evidence
    6. When an edge uses any other arrow
       - Then it fails as unsupported
    7. When the repo has no views, or a pull request changes only tests
       - Then it passes with no files; other runs read everything
