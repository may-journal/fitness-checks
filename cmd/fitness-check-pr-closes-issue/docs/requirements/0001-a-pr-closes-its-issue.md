# 0001 A PR Closes Its Issue

## Why

I see each merged PR close the issue it delivers.

## Measurement

PR bodies that close an issue

-

PR bodies

## Requirements

1. A PR body closes its issue
   1. GIVEN Developer has a PR body
      1. WHEN it closes an issue with any GitHub closing keyword and any reference form
         1. THEN `pr-closes-issue` passes it
      2. WHEN it only references an issue or spaces the number
         1. THEN it fails for lack of a closing keyword
      3. WHEN it says it implements an issue without closing it
         1. THEN it fails naming each unclosed issue
      4. WHEN the keyword sits only inside code or an HTML comment
         1. THEN it fails for lack of a closing keyword
      5. WHEN `--require-close` names an issue the body leaves open
         1. THEN it fails naming that issue
      6. WHEN the body closes every required issue
         1. THEN it passes
2. A missing PR body passes
   1. GIVEN Developer runs a check with no PR body
      1. WHEN `pr-closes-issue` runs
         1. THEN it passes
         2. THEN it judges zero files
