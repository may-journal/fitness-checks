# 0004 Writing Stays Clear

## Why

I can read any doc, issue, PR or comment in a minute, because it is short, plain and spelled right. My code comments hold no dead code or filler.

## Measurement

docs, bodies and comments that pass every writing check

-

docs, bodies and comments the writing checks judge

## Requirements

- 0004.1
  - Given Developer has prose with a long sentence, paragraph, section, list or dense passage
    - When the prose limit checks read it
      - Then each fails naming place and limit, and plain prose and step chains pass
- 0004.2
  - Given Developer has text with a misspelling, a forbidden phrase, bold, italic or a reframed sentence
    - When the spelling and wording checks read it
      - Then they fail naming each finding, and pass known terms, code and quotes
- 0004.3
  - Given Developer has code comments with long prose, dead code, filler or tutorial tone
    - When the comment checks read the repo
      - Then they fail naming each comment and rule, and pass why-comments
- 0004.4
  - Given Developer has Markdown files with wrong file name case or links to missing files
    - When the file name and link checks read the repo
      - Then they fail naming each file, and pass the repo's own case and real links
- 0004.5
  - Given Developer has Markdown docs and ADR files that differ from their folder's templates
    - When the doc template checks read the repo
      - Then they fail naming each difference, and pass docs that match and ADRs kept as issues
- 0004.6
  - Given Developer has an issue or Plan body with unchecked items or a broken Plan template
    - When the issue and Plan checks read the body
      - Then they fail naming the item or section, and pass bodies that follow the template
- 0004.7
  - Given Developer has a PR body that breaks its template or leaves a delivered issue open
    - When the PR body checks read it
      - Then they fail naming the section or issue, and pass bodies that follow it
- 0004.8
  - Given Developer has a comment that runs long, adds noise or keeps a list of open questions
    - When the comment and question checks read it
      - Then they fail naming the line, and pass short comments and questions asked in chat
