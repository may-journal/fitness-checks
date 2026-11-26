# 0004 Writing Stays Clear

## Why

I can read any doc, issue, PR or comment in a minute, because it is short, plain and spelled right. My code comments hold no dead code or filler.

## Measurement

docs, bodies and comments that pass every writing check

-

docs, bodies and comments the writing checks judge

## Requirements

1. Prose reads easily
   1. GIVEN Developer has a dense passage
      1. WHEN the readability check reads it
         1. THEN it fails naming place and limit
         2. THEN plain prose passes
2. Wording is plain
   1. GIVEN Developer has text with bold
      1. WHEN the wording checks read it
         1. THEN each fails naming the finding
         2. THEN code and quotes pass
3. Comments say why
   1. GIVEN Developer has code comments that are hard to read
   2. GIVEN they reframe a sentence
      1. WHEN the comment checks read the repo
         1. THEN each fails naming the comment and rule
         2. THEN why comments pass
4. Removed in 2026.11.26.0200: moved to the markdown-filename-camel-case and markdown-filename-kebab-case and markdown-links docs
5. Removed in 2026.11.26.0100: moved to the `doc-template` and `adr-as-issue` docs
6. Removed in 2026.11.26.0200: moved to the issue-checklist doc
7. Removed in 2026.11.26.0000: moved to the pr-closes-issue and pr-structure checks' own requirements
8. Comments stay short
   1. GIVEN Developer has a comment that lists open questions
      1. WHEN the question check reads it
         1. THEN it fails naming the line
         2. THEN questions asked in chat pass
