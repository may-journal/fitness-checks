# 0001 Givens Start With Declared Actors

## Why

I keep every role a requirement names one my system context declares, so a dropped role fails each Given naming it.

## Measurement

Givens that start with a declared actor

-

Givens in repos with architecture views

## Requirements

```gherkin
SCENARIO: 0001.1
  GIVEN Developer has a repo with requirement docs AND a system context declaring its actors
  WHEN a Given starts with an actor's name or Mermaid ID as a whole word in the diagram's case
  THEN it passes counting the docs read
    BUT some actors may start no Given
  WHEN a Given starts with no actor
    AND a Given starts with a plural
    AND a Given starts with another case
    AND a Given starts with an undeclared name
  THEN it fails naming the Given line
    AND it names the context line
    AND it names the declared actors
  WHEN several actors start a Given joined by `and`
  THEN it passes
    BUT it fails at a word after `and` that names no actor
  WHEN one actor's ID spells another actor's name
  THEN it fails as ambiguous naming the Given line and both context lines
  WHEN a Given names a system
    AND a Given names a legend
    AND a Given names a flowchart `person` node
    AND a Given names another view's person
  THEN it fails
    BUT only the context's C4 `Person` and `Person_Ext` are actors
  WHEN the repo has no architecture views
  THEN it passes with no files
  WHEN a doc's shape is wrong
    AND a view other than the context is missing
  THEN it reports neither
    AND it leaves them to `requirements` and `architecture-contract`

SCENARIO: 0001.2
  GIVEN Developer has a Given naming an actor the system context no longer declares
  WHEN a pull request changes only tests or changelogs
  THEN it passes with no files
    BUT such changes cannot move the views
  WHEN any other run changes even an unrelated file
  THEN it reads every doc
    AND it fails each Given naming that actor
```
