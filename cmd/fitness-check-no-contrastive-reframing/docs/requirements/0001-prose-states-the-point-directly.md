# 0001 Prose States The Point Directly

## Why

I read prose that says what a thing is. It never rejects a claim to restate it.

## Measurement

Markdown files and code comments with no not-X-it's-Y reframing

-

Markdown files and code comment blocks

## Requirements

1. Prose never reframes
   1. GIVEN Developer has Markdown prose or code comments
   2. GIVEN some hold a reframing sentence
      1. WHEN a negation is followed by a restating sentence
         1. THEN it fails once per match
         2. THEN it quotes the snippet
      2. WHEN one sentence negates plus pivots with a contrast
         1. THEN it fails quoting the snippet
      3. WHEN a negation has no restatement or the pattern sits in front matter or fenced code or a heading
         1. THEN it passes
         2. THEN the form `not only X but also Y` passes
         3. THEN the same holds in inline code or a quoted span
      4. WHEN the matching sentence is longer than 80 characters
         1. THEN the quoted snippet is shortened with an ellipsis
      5. WHEN a code comment holds the pattern
         1. THEN it fails naming the file
         2. THEN it names the line of the block
