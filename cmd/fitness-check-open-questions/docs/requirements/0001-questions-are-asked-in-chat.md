# 0001 Questions Are Asked In Chat

## Why

I get open questions in chat and never as a list buried in a doc, page or body.

## Measurement

docs, pages and bodies with no list of open questions

-

docs, pages and bodies

## Requirements

1. Lists of open questions fail
   1. GIVEN Developer has docs and pages
   2. GIVEN bodies may list open questions
      1. WHEN a heading or a bold line names open questions
         1. THEN `open-questions` fails naming the file and line
         2. THEN it gives the fix to ask in chat
      2. WHEN a label line names open questions
         1. THEN it fails naming the file and line
      3. WHEN open questions appear only in a sentence
         1. THEN it passes
      4. WHEN open questions appear only in code or a comment
         1. THEN it passes
      5. WHEN an HTML heading from `<h1>` to `<h6>` names open questions
         1. THEN it fails naming the file and line
      6. WHEN a body read with `--body-file` holds a label line of open questions
         1. THEN it fails naming the body
