# 0001 Comments Stay Short And Clean

## Why

I read short issue and PR comments that hold no agent noise, so each reply is worth my time.

## Measurement

comments that stay within the word and table budgets and hold no agent noise

-

comments checked

## Requirements

1. Comments stay within their budgets
   1. GIVEN Developer has a comment that reaches the check by `--body-file`
      1. WHEN the comment holds more words outside tables or code or diagrams or more plain tables than its budget
         1. THEN it fails naming the count
         2. THEN it fails naming the budget
      2. WHEN the repo sets `comment.maxWords` above the word count
         1. THEN it passes
2. Diagrams carry numbered callouts
   1. GIVEN Developer has a comment with a mermaid diagram
      1. WHEN the diagram has no numbered callouts
         1. THEN it fails naming the line
      2. WHEN the diagram is the legend right after a numbered diagram
         1. THEN it passes
3. Agent noise fails
   1. GIVEN Developer has a comment written by an agent
      1. WHEN the comment holds agent noise outside code such as a change-log note or agent frame or model identifier or fitness marker
         1. THEN it fails naming the line
         2. THEN it fails naming the noise
      2. WHEN the comment names a model inside a code span
         1. THEN it passes
4. Other input is left alone
   1. GIVEN Developer has text that is no comment
      1. WHEN the text is a Plan or Epic body piped under `--blob`
         1. THEN the check is skipped
      2. WHEN no comment reaches the check
         1. THEN it passes with no files
