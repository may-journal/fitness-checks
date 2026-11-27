# 0001 Docs Hold No Fitness Front Matter

## Why

I copy docs as they stand. A leftover `fitnessFunctions` or `relatedConfigurations` block that no check reads never spreads.

## Measurement

Markdown files holding no fitness front matter key

-

Markdown files

## Requirements

1. Fitness front matter fails
   1. GIVEN Developer has Markdown files
   2. GIVEN some hold fitness front matter
      1. WHEN a file or a PR template comment holds `fitnessFunctions` or `relatedConfigurations` in its leading block
         1. THEN it fails once per file
         2. THEN the failure names `fitness strip-front-matter`
      2. WHEN the file has CRLF line endings or a byte order mark
         1. THEN it still fails
         2. THEN other keys beside the fitness key do not save it
      3. WHEN the front matter holds only other keys
         1. THEN the file passes
      4. WHEN a leading rule sets a fitness key between two rules before prose
         1. THEN the file passes with the rule read as prose
      5. WHEN the file sits under `testdata` or `tests/examples`
         1. THEN it is never judged
