# 0001 Code View Matches The Source

## Why

I keep my code view's classes and edges true to the declarations and relations my source holds.

## Measurement

code view classes and edges that match the source

-

code view classes and edges

## Requirements

```gherkin
SCENARIO: 0001.1
  GIVEN Developer has a code view whose concrete classes link their declaring sources
  WHEN each label AND link AND kind AND signature AND member matches
  THEN it passes
    AND it allows generics
    AND it allows extensions
    AND it allows any spacing
    AND it allows methods elsewhere
  WHEN a class's linked lines lack its symbol
  THEN it fails naming where it is declared if anywhere
  WHEN a kind differs
    AND an interface mark differs
    AND a signature differs
    AND a member differs
    AND a field type differs
  THEN it fails naming the class line and the source line
  WHEN one name is declared twice under build conditions or as overloads
  THEN it fails listing both
    BUT a link or signature that picks one ends the failure
  WHEN a class links a test
    AND a class links a generated file
    AND a class links a non-source file
  THEN it fails
    AND a non-source link is unsupported evidence
  WHEN two classes name one declaration
    AND a label names no symbol
  THEN it fails
    BUT a class marked `<<abstract>>` is spared
  WHEN a shown folder holds an exported type no class names
  THEN it fails
    BUT functions and fields may be omitted

SCENARIO: 0001.2
  GIVEN Developer has a code view drawing edges between its classes
  WHEN the source uses a class no edge allows even via helpers
  THEN it fails at the source line naming the smallest holding class
  WHEN the edge is `-->` or `..>` or `*--` or `o--` or `--` or `..`
  THEN it allows uses from the first class
    BUT the undirected `--` and `..` allow uses either way
  WHEN a drawn association has no source relation
    AND a drawn realization has no source relation
    AND a drawn inheritance has no source relation
  THEN it fails at the edge
  WHEN an edge row links a source
  THEN the linked lines must hold that relation
  WHEN an edge touches an abstract class
  THEN it needs no use
    BUT a row link is unsupported evidence
  WHEN an edge uses any other arrow
  THEN it fails as unsupported
  WHEN the repo has no views
    AND a pull request changes only tests
  THEN it passes with no files
    BUT other runs read everything
```
