# 0004 Writing Stays Clear

## Why

I can read any doc, issue, PR or comment in a minute, because it is short, plain and spelled right. My code comments hold no dead code or filler.

## Measurement

docs, bodies and comments that pass every writing check

-

docs, bodies and comments the writing checks judge

## Requirements

1. Prose stays within limits
   1. GIVEN Developer has prose with a long sentence
   2. GIVEN it has a long section
   3. GIVEN it has a long list
   4. GIVEN it has a dense passage
      1. WHEN the prose limit checks read it
         1. THEN each fails naming place and limit
         2. THEN plain prose and step chains pass
2. Wording is plain
   1. GIVEN Developer has text with a misspelling
   2. GIVEN it has a forbidden phrase
   3. GIVEN it has bold
      1. WHEN the wording checks read it
         1. THEN each fails naming the finding
         2. THEN known terms and code and quotes pass
3. Comments say why
   1. GIVEN Developer has code comments with long prose or dead code
      1. WHEN the comment checks read the repo
         1. THEN each fails naming the comment and rule
         2. THEN why comments pass
4. Markdown links work
   1. GIVEN Developer has Markdown files with wrong name case or dead links
      1. WHEN the file name and link checks read the repo
         1. THEN each fails naming the file
         2. THEN real links pass
5. Docs follow templates
   1. GIVEN Developer has docs and ADRs that differ from their templates
      1. WHEN the doc template checks read the repo
         1. THEN each difference fails
         2. THEN matching docs pass
6. Issues follow templates
   1. GIVEN Developer has an issue or Plan body that breaks its template
      1. WHEN the issue and Plan checks read the body
         1. THEN each fails naming the item or section
         2. THEN bodies that follow the template pass
7. PR bodies follow templates
   1. GIVEN Developer has a PR body that breaks its template
      1. WHEN the PR body checks read it
         1. THEN each fails naming the section
         2. THEN bodies that follow it pass
8. Comments stay short
   1. GIVEN Developer has a comment that runs long
   2. GIVEN it adds noise
   3. GIVEN it lists open questions
      1. WHEN the comment and question checks read it
         1. THEN each fails naming the line
         2. THEN short comments pass
