# 0004 Writing Stays Clear

## Why

I can read any doc, issue, PR or comment in a minute, because it is short, plain and spelled right. My code comments hold no dead code or filler.

## Measurement

docs, bodies and comments that pass every writing check

-

docs, bodies and comments the writing checks judge

## Requirements

```gherkin
SCENARIO: 0004.1
  GIVEN Developer has prose with a long sentence or paragraph or section or list or dense passage
  WHEN the prose limit checks read it
  THEN each fails naming place and limit
    AND plain prose and step chains pass

SCENARIO: 0004.2
  GIVEN Developer has text with bold or italic or a reframed sentence
  WHEN the wording checks read it
  THEN they fail naming each finding
    AND they pass code and quotes

SCENARIO: 0004.3
  GIVEN Developer has code comments with long prose or dead code or filler or tutorial tone
  WHEN the comment checks read the repo
  THEN they fail naming each comment and rule
    AND they pass why-comments

SCENARIO: 0004.4
  GIVEN Developer has Markdown files with wrong file name case or links to missing files
  WHEN the file name and link checks read the repo
  THEN they fail naming each file
    AND they pass the repo's own case and real links

SCENARIO: 0004.5
  GIVEN Developer has Markdown docs and ADR files that differ from their folder's templates
  WHEN the doc template checks read the repo
  THEN they fail naming each difference
    AND they pass docs that match and ADRs kept as issues

SCENARIO: 0004.6
  GIVEN Developer has an issue or Plan body with unchecked items or a broken Plan template
  WHEN the issue and Plan checks read the body
  THEN they fail naming the item or section
    AND they pass bodies that follow the template

SCENARIO: 0004.7
  GIVEN Developer has a PR body that breaks its template or leaves a delivered issue open
  WHEN the PR body checks read it
  THEN they fail naming the section or issue
    AND they pass bodies that follow it

SCENARIO: 0004.8
  GIVEN Developer has a comment that runs long or adds noise or keeps a list of open questions
  WHEN the comment and question checks read it
  THEN they fail naming the line
    AND they pass short comments and questions asked in chat
```
