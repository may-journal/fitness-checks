# 0001 Markdown Holds No Bold Or Italic

## Why

I read plain Markdown in which no bold or italic emphasis shouts or fades.

## Measurement

Markdown files with no bold or italic emphasis

-

Markdown files

## Requirements

1. Emphasis fails by occurrence
   1. GIVEN Developer has Markdown files
   2. GIVEN some hold bold or italic emphasis
      1. WHEN a file holds asterisk or underscore bold or italic
         1. THEN it fails once per occurrence
         2. THEN each failure names the file
         3. THEN each failure names the snippet
      2. WHEN the emphasis sits in fenced code or inline code
         1. THEN the file passes
         2. THEN emphasis in a link block passes the same way
      3. WHEN a line starts with an asterisk list marker
         1. THEN that line passes
      4. WHEN Developer gives a body with `--body-file`
         1. THEN the body is judged as the description
      5. WHEN `CHANGELOG.md` holds bold outside a Release Please scope
         1. THEN only that bold fails
      6. WHEN the repo has no Markdown files
         1. THEN the check passes with no files
