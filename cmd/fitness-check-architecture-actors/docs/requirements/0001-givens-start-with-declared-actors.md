# 0001 Givens Start With Declared Actors

## Why

I keep every role a requirement names one my system context declares, so a dropped role fails each Given naming it.

## Measurement

Givens that start with a declared actor

-

Givens in repos with architecture views

## Requirements

1. Givens start with declared actors
   1. GIVEN Developer has a repo with requirement docs
   2. GIVEN it has a system context declaring its actors
      1. WHEN an actor's name or Mermaid ID opens a GIVEN as a whole word in the diagram's case
         1. THEN it passes counting the docs read
         2. THEN it passes though some actors start no GIVEN
      2. WHEN no declared actor opens a GIVEN
         1. THEN it fails naming the GIVEN line
         2. THEN it names the context line
         3. THEN it names the declared actors
         4. THEN a plural or another case or an undeclared name fails the same way
      3. WHEN several actors start a GIVEN joined by `and`
         1. THEN it passes
         2. THEN it fails at a word after `and` that names no actor
      4. WHEN one actor's ID spells another actor's name
         1. THEN it fails as ambiguous naming the GIVEN line
         2. THEN it names both context lines
      5. WHEN a GIVEN names something other than an actor
         1. THEN a system or a legend fails
         2. THEN a flowchart `person` node fails
         3. THEN another view's person fails
         4. THEN only the context's C4 `Person` or `Person_Ext` are actors
      6. WHEN the repo has no architecture views
         1. THEN it passes with no files
      7. WHEN a doc's shape is wrong or it holds bullet acceptances or a view other than the context is missing
         1. THEN it reports none of them
         2. THEN it leaves them to `requirements`
         3. THEN it leaves them to `architecture-contract`
2. A dropped actor fails every GIVEN naming it
   1. GIVEN Developer has a GIVEN naming an actor the system context no longer declares
      1. WHEN a pull request changes only tests or changelogs
         1. THEN it passes with no files
         2. THEN such changes cannot move the views
      2. WHEN any other run changes even an unrelated file
         1. THEN it reads every doc
         2. THEN it fails each GIVEN naming that actor
