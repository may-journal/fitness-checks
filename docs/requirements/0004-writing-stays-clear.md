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
4. Markdown links work
   1. GIVEN Developer has Markdown files with wrong name case or dead links
      1. WHEN the file name and link checks read the repo
         1. THEN each fails naming the file
         2. THEN real links pass
5. Docs follow templates
   1. GIVEN Developer has ADR files in folders or kept as issues
      1. WHEN the ADR checks read the repo
         1. THEN each ADR file in a folder fails
         2. THEN ADRs kept as issues pass
6. Issues are ticked
   1. GIVEN Developer has an issue body with unchecked items
      1. WHEN the issue checklist check reads the body
         1. THEN it fails naming the item
         2. THEN bodies with every item ticked pass
7. Removed in 2026.11.26.0000: moved to the pr-closes-issue and pr-structure checks' own requirements
8. Comments stay short
   1. GIVEN Developer has a comment that lists open questions
      1. WHEN the question check reads it
         1. THEN it fails naming the line
         2. THEN questions asked in chat pass
