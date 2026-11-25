# 0001 Givens Start With Declared Actors

## Why

I keep every role a requirement names one my system context declares, so a dropped role fails each Given naming it.

## Measurement

Givens that start with a declared actor

-

Givens in repos with architecture views

## Requirements

- 0001.1
  - Given Developer has a repo with requirement docs and a system context declaring its actors
    1. When a Given starts with an actor's name or Mermaid ID, as a whole word in the diagram's case
       - Then it passes, counting the docs read, though some actors start no Given
    2. When a Given starts with no actor, a plural, another case, or an undeclared name
       - Then it fails naming the Given line, the context line and the declared actors
    3. When several actors start a Given joined by `and`
       - Then it passes, but fails at a word after `and` that names no actor
    4. When one actor's ID spells another actor's name
       - Then it fails as ambiguous, naming the Given line and both context lines
    5. When a Given names a system, a legend, a flowchart `person` node, or another view's person
       - Then it fails, since only the context's C4 `Person` and `Person_Ext` are actors
    6. When the repo has no architecture views
       - Then it passes with no files
    7. When a doc's shape is wrong, or a view other than the context is missing
       - Then it reports neither, leaving them to `requirements` and `architecture-contract`
- 0001.2
  - Given Developer has a Given naming an actor the system context no longer declares
    1. When a pull request changes only tests or changelogs
       - Then it passes with no files, since such changes cannot move the views
    2. When any other run changes even an unrelated file
       - Then it reads every doc and fails each Given naming that actor
