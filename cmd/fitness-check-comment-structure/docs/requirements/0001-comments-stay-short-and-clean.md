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
      1. WHEN the comment holds more words than the budget outside tables and code and diagrams
         1. THEN it fails naming the count and the budget
      2. WHEN the repo sets `comment.maxWords` above that count
         1. THEN it passes
      3. WHEN the comment holds more plain tables than the budget
         1. THEN it fails naming the count and the budget
2. Diagrams carry numbered callouts
   1. GIVEN Developer has a comment with a mermaid diagram
      1. WHEN the diagram has no numbered callouts
         1. THEN it fails naming the line
      2. WHEN the diagram is the legend right after a numbered diagram
         1. THEN it passes
3. Agent noise fails
   1. GIVEN Developer has a comment written by an agent
      1. WHEN the comment holds a change-log note such as `Edit 2:`
         1. THEN it fails naming the line and the note
      2. WHEN the comment holds an agent frame such as a `<system-reminder>` tag
         1. THEN it fails naming the line and the frame
      3. WHEN the comment holds a model identifier outside code
         1. THEN it fails naming the line and the identifier
      4. WHEN the comment holds a fitness marker outside code
         1. THEN it fails naming the line and the marker
      5. WHEN the comment names a model inside a code span
         1. THEN it passes
4. Other input is left alone
   1. GIVEN Developer has text that is no comment
      1. WHEN the text is a Plan or Epic body piped under `--blob`
         1. THEN the check is skipped
      2. WHEN no comment reaches the check
         1. THEN it passes with no files
